package centrifugeinput

import (
	"math"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	utils "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/utils"
)

const (
	// Output
	InputName          = "centrifuge_input"
	FieldName_endpoint = "endpoint"
	FieldName_channel  = "channel"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an input of a centrifuge channel.").
	Field(benthos_service.NewStringField(FieldName_endpoint)).
	Field(benthos_service.NewStringField(FieldName_channel))

const useGlobalSkipFrameCount = math.MinInt32

func (s *service) Register() error {
	benthos_service.RegisterInput(InputName,
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Input, err error) {
			s.logger = mgr.Logger()
			log := s.logger.With("input", InputName).With("a", 1)
			log.With("caller", utils.Caller()).Info("Register")

			endpoint, err := conf.FieldString(FieldName_endpoint)
			if err != nil {
				return nil, err
			}
			s.endpoint = endpoint

			channel, err := conf.FieldString(FieldName_channel)
			if err != nil {
				return nil, err
			}
			s.channel = channel

			return s, nil
		})
	return nil
}
