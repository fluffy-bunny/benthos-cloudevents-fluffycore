package main

var jsonSchema = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","properties":{"id":{"type":"string"},"source":{"type":"string"},"specVersion":{"type":"string","enum":["1.0"]},"type":{"type":"string","enum":["requestunits.v1"]},"attributes":{"type":"object","properties":{"orgid":{"type":"object","properties":{"ceString":{"type":"string"}},"required":["ceString"]},"time":{"type":"object","properties":{"ceTimestamp":{"type":"string"}},"required":["ceTimestamp"]}},"required":["orgid","partition-key","time"]},"textData":{"type":"string"}},"required":["id","source","specVersion","type","attributes","textData"]}`

var kafkaYaml = `
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
            schema: "${JSON_SCHEMA}"
        #            schema_path: "file://C:/work/mapped/benthos-cloudevents-fluffycore/config/benthos/request_units_schema.json"
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
                          - mapping: "deadletterit(@,content())"
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
    "@service": benthos

`
