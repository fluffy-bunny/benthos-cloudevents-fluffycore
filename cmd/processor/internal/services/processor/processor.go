package processor

import (
	"context"
	"encoding/json"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
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
	generics := []map[string]interface{}{}
	for _, item := range request.BadBatch.Messages {
		generic := make(map[string]interface{})
		err := json.Unmarshal(item, &generic)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal generic")
			continue
		}
		generics = append(generics, generic)
	}
	log.Info().Interface("generics", generics).Msg("processBadBatch")

	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil

}
