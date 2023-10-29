package input

import (
	"context"
	"errors"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/service"
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

var gibberishConfigSpec = service.NewConfigSpec().
	Summary("Creates an input that generates garbage.").
	Field(service.NewIntField("length").Default(100))

func newGibberishInput(conf *service.ParsedConfig) (service.Input, error) {
	length, err := conf.FieldInt("length")
	if err != nil {
		return nil, err
	}
	if length <= 0 {
		return nil, fmt.Errorf("length must be greater than 0, got: %v", length)
	}
	if length > 10000 {
		return nil, errors.New("that length is way too high bruh")

	}
	return service.AutoRetryNacks(&gibberishInput{length}), nil
}

func init() {
	err := service.RegisterInput(
		"gibberish", gibberishConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Input, error) {
			return newGibberishInput(conf)
		})
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type gibberishInput struct {
	length int
}

func (g *gibberishInput) Connect(ctx context.Context) error {
	return nil
}

func (g *gibberishInput) Read(ctx context.Context) (*service.Message, service.AckFunc, error) {
	//log := zerolog.Ctx(ctx).With().Str("input", "gibberish").Logger()
	ceBatch := &cloudevents.CloudEventBatch{}
	for i := 0; i < g.length; i++ {
		ce := &cloudevents.CloudEvent{
			Id:          uuid.New().String(),
			Source:      "https://gibberish.com/cloudevent",
			Type:        "CloudEvent",
			SpecVersion: "1.0",
			Data: &cloudevents.CloudEvent_TextData{
				TextData: `{"a":"b"}`,
			},
		}
		ceBatch.Events = append(ceBatch.Events, ce)
	}

	b, err := protojson.Marshal(ceBatch)
	if err != nil {
		return nil, nil, err
	}

	return service.NewMessage(b), func(ctx context.Context, err error) error {

		// Nacks are retried automatically when we use service.AutoRetryNacks
		return nil
	}, nil
}

func (g *gibberishInput) Close(ctx context.Context) error {
	return nil
}
