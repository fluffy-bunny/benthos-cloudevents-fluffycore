package centrifugestream

import (
	"context"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	service struct {
		contracts_benthos.UnimplementedIBenthosStream
		config *contracts_config.Config
		stream *benthos_service.Stream
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_benthos.IBenthosStream = (*service)(nil)
}

func AddSingletonIBenthosStream(builder di.ContainerBuilder) {
	di.AddSingleton[contracts_benthos.IBenthosStream](builder, stemService.Ctor)
}

func (s *service) Ctor(config *contracts_config.Config) (contracts_benthos.IBenthosStream, error) {
	builderOne := benthos_service.NewStreamBuilder()
	err := builderOne.SetYAML(benthosConfig)
	if err != nil {
		return nil, err
	}
	stream, err := builderOne.Build()
	if err != nil {
		return nil, err
	}

	return &service{
		config: config,
		stream: stream,
	}, nil

}

func (s *service) Run(ctx context.Context) (err error) {
	return s.stream.Run(ctx)
}
func (s *service) Stop(ctx context.Context) (err error) {
	return s.stream.Stop(ctx)
}

const benthosConfig = `
http:
  enabled: false
input:
  centrifuge_input:
    channel: chat:index
    batching:
      count: 3
      period: 20s
pipeline:
  threads: 1
  processors:
    - sleep:
        duration: 1s
output:
  justlogit:
    max_in_flight: 64
logger:
  level: ${LOG_LEVEL}
  format: json
  add_timestamp: true
  static_fields:
    "@service": benthos.kafka

`
