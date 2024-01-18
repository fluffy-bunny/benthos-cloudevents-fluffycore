package kafkastream

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
    topics: ["cloudevents-core"]
    consumer_group: "$Default2"
    multi_header: true
    batching:
      count: 3
      period: 20s
      processors:
        - json_schema:
            schema: '{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","properties":{"id":{"type":"string"},"source":{"type":"string"},"specVersion":{"type":"string","enum":["1.0"]},"type":{"type":"string","enum":["requestunits.v1"]},"attributes":{"type":"object","properties":{"orgid":{"type":"object","properties":{"ceString":{"type":"string"}},"required":["ceString"]},"time":{"type":"object","properties":{"ceTimestamp":{"type":"string"}},"required":["ceTimestamp"]}},"required":["orgid","partition-key","time"]},"textData":{"type":"string"}},"required":["id","source","specVersion","type","attributes","textData"]}'
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
  cloudevents_grpc:
    grpc_url: "${OUTPUT_CLOUDEVENTOUTPUT_GRPC_URL}"
    max_in_flight: 64
    channel: "mychannel"

    # auth[optional] one of: [oauth2,basic,apikey](ordered by priority)
    #--------------------------------------------------------------------
    #auth:
    #  basic:
    #    user_name: "admin"
    #    password: "password"
    #  oauth2:
    #    client_id: "my_client_id"
    #    client_secret: "secret"
    #    token_endpoint: "https://example.com/oauth2/token"
    #    scopes: ["scope1", "scope2"]
    #  apikey:
    #    name: "x-api-key"
    #    value: "secret"
logger:
  level: ${LOG_LEVEL}
  format: json
  add_timestamp: true
  static_fields:
    "@service": benthos.kafka
`
