package kafkacloudeventservice

import (
	"context"
	"sync"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	kafka_franz "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/kafka_franz/v2"
	proto_kafkacloudevent "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/kafkacloudevent"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

type (
	service struct {
		proto_kafkacloudevent.UnimplementedKafkaCloudEventServiceServer
		config      *contracts_config.Config
		KafkaClient *kgo.Client
	}
)

func AddKafkaCloudEventServiceServer(builder di.ContainerBuilder) {
	proto_kafkacloudevent.AddKafkaCloudEventServiceServer[proto_kafkacloudevent.IKafkaCloudEventServiceServer](builder,
		func(config *contracts_config.Config, kafkaClient *kgo.Client) (proto_kafkacloudevent.IKafkaCloudEventServiceServer, error) {
			return &service{
				config:      config,
				KafkaClient: kafkaClient,
			}, nil
		})
}

var prefix = "ce_"

func (s *service) SubmitCloudEvents(ctx context.Context, request *proto_kafkacloudevent.SubmitCloudEventsRequest) (*proto_kafkacloudevent.SubmitCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Interface("request", request).Msg("request")

	for _, event := range request.Batch.Events {
		record, err := kafka_franz.CloudEventToKafkaMessage(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert CloudEvent to KafkaMessage")
			return nil, err
		}
		record.Topic = s.config.KafkaConfig.Topic

		var wg sync.WaitGroup
		wg.Add(1)

		s.KafkaClient.Produce(ctx, record, func(_ *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				log.Error().Err(err).Msg("failed to produce")
			}

		})
		wg.Wait()
	}

	return &proto_kafkacloudevent.SubmitCloudEventsResponse{}, nil
}
