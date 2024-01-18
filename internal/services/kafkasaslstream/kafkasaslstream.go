package kafkasaslstream

import (
	"context"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	service struct {
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
  kafka:
    addresses: ["${INPUT_KAFKA_BROKERS}"]
    topics: ["mapped-api-usage"]
    consumer_group: "$Default2"
    tls:
      enabled: true
    sasl:
      mechanism: PLAIN
      user: "${INPUT_SASL_USERNAME}"
      password: "${INPUT_SASL_PASSWORD}"
    multi_header: true
    batching:
      count: 3
      period: 60s
      processors:
        - json_schema:
            schema: '{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","properties":{"requestTime":{"type":"string"},"requestUnits":{"type":"integer"},"orgId":{"type":"string"},"tokenId":{"type":"string"},"uniquenessToken":{"type":"string"}},"required":["requestTime","requestUnits","orgId","tokenId","uniquenessToken"]}'
        - switch:
            - check: errored()
              processors:
                - for_each:
                    - while:
                        at_least_once: true
                        max_loops: 0
                        check: errored()
                        processors:
                          - catch: [] # Wipe any previous error
                          - mapping: "errorlogit(@,content())"
                - mapping: |
                    deleted()

        - archive:
            format: binary

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
    "@service": benthos.kakfa.sasl
`
