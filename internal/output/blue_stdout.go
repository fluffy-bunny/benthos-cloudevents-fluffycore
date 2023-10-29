package output

import (
	"context"

	"github.com/benthosdev/benthos/v4/public/service"
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	err := service.RegisterOutput(
		"blue_stdout", service.NewConfigSpec(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (out service.Output, maxInFlight int, err error) {
			return &blueOutput{}, 1, nil
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type blueOutput struct{}

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
