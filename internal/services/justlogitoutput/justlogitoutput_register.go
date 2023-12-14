package justlogitoutput

import (
	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	utils "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/utils"
)

const (
	// Output
	OutputName            = "justlogit"
	FieldName_maxInFlight = "max_in_flight"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an output to a grpc service.").
	Field(benthos_service.NewInputMaxInFlightField())

func (s *service) Register() error {
	benthos_service.RegisterOutput(OutputName,
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Output, maxInFlight int, err error) {
			s.logger = mgr.Logger()
			log := s.logger.With("output", OutputName).With("a", 1)
			log.With("caller", utils.Caller()).Info("Register")
			log.Debug("hi")

			maxInFlight, err = conf.FieldInt(FieldName_maxInFlight)
			if err != nil {
				return nil, 0, err
			}
			return s, maxInFlight, nil
		})
	return nil
}
