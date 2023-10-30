package kafkacloudeventservice

import (
	"context"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	proto_kafkacloudevent "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/kafkacloudevent"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	"github.com/rs/zerolog"
)

type (
	service struct {
		proto_kafkacloudevent.UnimplementedKafkaCloudEventServiceServer
		config *contracts_config.Config
	}
)

func AddKafkaCloudEventServiceServer(builder di.ContainerBuilder) {
	proto_kafkacloudevent.AddKafkaCloudEventServiceServer[proto_kafkacloudevent.IKafkaCloudEventServiceServer](builder,
		func(config *contracts_config.Config) proto_kafkacloudevent.IKafkaCloudEventServiceServer {
			return &service{
				config: config,
			}
		})
}

func (s *service) SubmitCloudEvents(ctx context.Context, request *proto_kafkacloudevent.SubmitCloudEventsRequest) (*proto_kafkacloudevent.SubmitCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Interface("request", request).Msg("request")
	return &proto_kafkacloudevent.SubmitCloudEventsResponse{}, nil
}
