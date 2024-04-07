package centrifugeinput

import (
	"math"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	utils "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/utils"
)

const (
	// Output
	InputName = "centrifuge_input"
	//FieldName_endpoint = "endpoint"
	FieldName_channel = "channel"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an input of a centrifuge channel.").
	Field(benthos_service.NewStringField(FieldName_channel)).
	Field(benthos_service.NewBatchPolicyField("batching").
		Description("Allows you to configure a [batching policy](/docs/configuration/batching) that applies to individual topic partitions in order to batch messages together before flushing them for processing. Batching can be beneficial for performance as well as useful for windowed processing, and doing so this way preserves the ordering of topic partitions.").
		Advanced())

const useGlobalSkipFrameCount = math.MinInt32

func (s *service) Register() error {
	benthos_service.RegisterInput(InputName,
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Input, err error) {
			s.logger = mgr.Logger()
			log := s.logger.With("input", InputName).With("a", 1)
			log.With("caller", utils.Caller()).Info("Register")

			channel, err := conf.FieldString(FieldName_channel)
			if err != nil {
				return nil, err
			}
			s.channel = channel

			return s, nil
		})
	return nil
}
