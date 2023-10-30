package output

import (
	"context"

	"github.com/benthosdev/benthos/v4/public/service"
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	log "github.com/rs/zerolog/log"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

var gibberishConfigSpec = service.NewConfigSpec().
	Summary("Creates an output to a grpc service.").
	Field(service.NewStringField("grpc_url"))

func init() {

	err := service.RegisterOutput(
		"blue_stdout", gibberishConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (out service.Output, maxInFlight int, err error) {
			return newBlueOutput(conf)
		})
	if err != nil {
		panic(err)
	}

}

//------------------------------------------------------------------------------

type blueOutput struct {
	grpcUrl string
}

func newBlueOutput(conf *service.ParsedConfig) (out service.Output, maxInFlight int, err error) {
	grpcUrl, err := conf.FieldString("grpc_url")
	if err != nil {
		return nil, 0, err
	}

	return &blueOutput{
		grpcUrl: grpcUrl,
	}, 1, nil
}

func (b *blueOutput) Connect(ctx context.Context) error {
	return nil
}

func (b *blueOutput) Write(ctx context.Context, msg *service.Message) error {

	content, err := msg.AsBytes()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert message to bytes")
		return err
	}
	ceBatch := &cloudevents.CloudEventBatch{}
	err = protojson.Unmarshal(content, ceBatch)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CloudEventBatch")
		return err
	}
	log.Info().Interface("ceBatch", ceBatch).Msg("ceBatch")
	return nil
}

func (b *blueOutput) Close(ctx context.Context) error {
	return nil
}
