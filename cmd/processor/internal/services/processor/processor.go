package processor

import (
	"context"
	"sync"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	contracts_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/kafkaclient"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

type (
	service struct {
		proto_cloudeventprocessor.UnimplementedCloudEventProcessorServer
		config           *contracts_config.Config
		deadLetterClient contracts_kafkaclient.IDeadLetterClient
	}
)

func AddCloudEventProcessorServer(builder di.ContainerBuilder) {
	proto_cloudeventprocessor.AddCloudEventProcessorServer[proto_cloudeventprocessor.ICloudEventProcessorServer](builder,
		func(config *contracts_config.Config, deadLetterClient contracts_kafkaclient.IDeadLetterClient) proto_cloudeventprocessor.ICloudEventProcessorServer {
			return &service{
				config:           config,
				deadLetterClient: deadLetterClient,
			}
		})
}

func (s *service) ProcessCloudEvents(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	s.processGoodBatch(ctx, request)
	s.processBadBatch(ctx, request)
	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
}

func (s *service) processGoodBatch(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("--> processGoodBatch")
	log.Info().Interface("good_batch", request.Batch).Msg("processGoodBatch")
	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
}

func (s *service) processBadBatch(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("--> processBadBatch")
	log.Info().Interface("bad_batch", request.BadBatch).Msg("processBadBatch")
	if request.BadBatch == nil {
		return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
	}
	for _, item := range request.BadBatch.Messages {
		var wg sync.WaitGroup
		wg.Add(1)
		record := &kgo.Record{
			Value: item,
		}
		kafkaClient := s.deadLetterClient.GetClient()

		kafkaClient.Produce(ctx, record, func(_ *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				log.Error().Err(err).Msg("failed to produce")
			}
		})
		wg.Wait()
	}

	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil

}
