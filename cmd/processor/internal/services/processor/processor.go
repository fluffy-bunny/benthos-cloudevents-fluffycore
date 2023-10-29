package processor

import (
	"context"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	"github.com/rs/zerolog"
)

type (
	service struct {
		proto_cloudeventprocessor.UnimplementedCloudEventProcessorServer
		config *contracts_config.Config
	}
)

func AddCloudEventProcessorServer(builder di.ContainerBuilder) {
	proto_cloudeventprocessor.AddCloudEventProcessorServer[proto_cloudeventprocessor.ICloudEventProcessorServer](builder,
		func(config *contracts_config.Config) proto_cloudeventprocessor.ICloudEventProcessorServer {
			return &service{
				config: config,
			}
		})
}

func (s *service) ProcessCloudEvents(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("ProcessCloudEvents")
	log.Info().Interface("request", request).Msg("request")
	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
}
