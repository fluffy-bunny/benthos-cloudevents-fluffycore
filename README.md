# benthos-cloudevents-fluffycore

Benthos plugin that processes cloudevents and sends them downstream to a gprc handler.  

## Build the proto

[grpc.io](https://grpc.io/docs/languages/go/basics/)  
[Transcode Ref](https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/)  
[custom protoc plugin](https://rotemtam.com/2021/03/22/creating-a-protoc-plugin-to-gen-go-code/)  

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fluffy-bunny/fluffycore/protoc-gen-go-fluffycore-di/cmd/protoc-gen-go-fluffycore-di@latest
```

```bash
protoc --go_out=. --go_opt=paths=source_relative ./pkg/proto/cloudevents/cloudevents.proto  

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --go-fluffycore-di_out=.  --go-fluffycore-di_opt=paths=source_relative ./pkg/proto/cloudeventprocessor/cloudeventprocessor.proto 

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --go-fluffycore-di_out=.  --go-fluffycore-di_opt=paths=source_relative ./pkg/proto/kafkacloudevent/kafkacloudevent.proto 

```

## Grpc References

[auth examples](https://github.com/johanbrandhorst/grpc-auth-example)  


## Docker

```bash
 docker build --file .\build\Dockerfile.processor . --tag fluffybunny.benthos.processor

 docker build --file .\build\Dockerfile.benthos . --tag fluffybunny.benthos.benthos

 docker-compose up
```

## Services

### Kafka UI

[Kafka-ui](http://localhost:9090/)  

### KafkaCloudEventService - grpc

```bash
grpc://localhost:9050
```

This is service to submit a cloud-event to kafka.  

#### Request

```json
{
    "batch": {
        "events": [
            {
                "attributes": [
                    {
                        "value": {
                            "ce_string": "test"
                        },
                        "key": "testkey"
                    }
                ],
                "spec_version": "1.0",
               
                "text_data": "{\"a\":\"b\"}",
                
                "id": "1234",
                "type": "my.type",
                "source": "//my/source"
            }
        ]
    }
}
```

#### Response 

```json
{}
```

### CloudEventProcessor - grpc

```bash
grpc://localhost:9050
```

This service receives cloud-events as batches via out custom benthos output handler.

**IMPORTANT**: This is an all or nothing process.  It sends the messages in batches.  There are good and bad and the processor needs to decide what it wants to do with the bad messages.  In some cases the entire batch is bad at the 

### Viewer

I use docker desktop and look at the logs.  My downstream processor just logs what it got and returns a nil.  

