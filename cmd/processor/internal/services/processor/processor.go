package processor

import (
	"context"
	"sync"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	contracts_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/kafkaclient"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	"github.com/gogo/status"
	zerolog "github.com/rs/zerolog"
	kgo "github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc/codes"
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
func (s *service) validateProcessCloudEventsRequest(request *proto_cloudeventprocessor.ProcessCloudEventsRequest) error {
	if request == nil {
		return nil
	}
	if fluffycore_utils.IsEmptyOrNil(request.Channel) {
		return status.Error(codes.InvalidArgument, "request.Channel is empty")
	}

	return nil
}
func (s *service) ProcessCloudEvents(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx).With().Str("method", "ProcessCloudEvents").Logger()
	err := s.validateProcessCloudEventsRequest(request)
	if err != nil {
		log.Debug().Err(err).Msg("validateProcessCloudEventsRequest")
		return nil, err
	}
	log = log.With().Str("channel", request.Channel).Logger()
	ctx = log.WithContext(ctx)
	s.processGoodBatch(ctx, request)
	s.processBadBatch(ctx, request)
	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
}

func (s *service) processGoodBatch(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx).With().Int("count", len(request.Batch.Events)).Logger()
	log.Info().Msg("--> processGoodBatch")
	log.Info().Interface("good_batch", request.Batch).Msg("processGoodBatch")
	return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
}

func (s *service) processBadBatch(ctx context.Context, request *proto_cloudeventprocessor.ProcessCloudEventsRequest) (*proto_cloudeventprocessor.ProcessCloudEventsResponse, error) {
	if request.BadBatch == nil {
		return &proto_cloudeventprocessor.ProcessCloudEventsResponse{}, nil
	}
	log := zerolog.Ctx(ctx).With().Int("count", len(request.BadBatch.Messages)).Logger()
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
