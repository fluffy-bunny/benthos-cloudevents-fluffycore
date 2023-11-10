package kafkacloudeventservice

import (
	"context"
	"fmt"
	"sync"

	"regexp"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	kafka_franz "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/kafka_franz/v2"
	proto_cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	proto_kafkacloudevent "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/kafkacloudevent"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	zerolog "github.com/rs/zerolog"
	kgo "github.com/twmb/franz-go/pkg/kgo"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type (
	service struct {
		proto_kafkacloudevent.UnimplementedKafkaCloudEventServiceServer
		config      *contracts_config.Config
		KafkaClient *kgo.Client
	}
)

func AddKafkaCloudEventServiceServer(builder di.ContainerBuilder) {
	proto_kafkacloudevent.AddKafkaCloudEventServiceServer[proto_kafkacloudevent.IKafkaCloudEventServiceServer](builder,
		func(config *contracts_config.Config, kafkaClient *kgo.Client) (proto_kafkacloudevent.IKafkaCloudEventServiceServer, error) {
			return &service{
				config:      config,
				KafkaClient: kafkaClient,
			}, nil
		})
}

var _reservedAttributes = []string{
	"partition-key", // what key to use when publishing to kafka
}
var _reservedAttributesMap = map[string]bool{}
var regexCompiled *regexp.Regexp

// https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#context-attributes
const pattern = "^[a-z0-9]{1,20}$"

func init() {
	for _, v := range _reservedAttributes {
		_reservedAttributesMap[v] = true
	}
	regexCompiled = regexp.MustCompile(pattern)
}
func (s *service) validateCloudEvent(request *proto_kafkacloudevent.SubmitCloudEventsRequest) error {
	ceTimestamp := timestamppb.Now()

	if fluffycore_utils.IsEmptyOrNil(request) {
		return status.Error(codes.InvalidArgument, "request is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Batch) {
		return status.Error(codes.InvalidArgument, "request.Batch is empty")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Batch.Events) {
		return status.Error(codes.InvalidArgument, "request.Batch.Events is empty")
	}

	ensureEventTimestamp := func(event *proto_cloudevents.CloudEvent) {
		setTime := func() {
			event.Attributes["time"] = &proto_cloudevents.CloudEvent_CloudEventAttributeValue{
				Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeTimestamp{
					CeTimestamp: ceTimestamp,
				},
			}
		}
		attrib, ok := event.Attributes["time"]
		if !ok {
			setTime()
		} else {
			timeStamp := attrib.GetCeTimestamp()
			if timeStamp == nil {
				setTime()
			}
		}
	}
	validateAttributes := func(event *proto_cloudevents.CloudEvent) error {
		for attribute := range event.Attributes {
			if fluffycore_utils.IsEmptyOrNil(attribute) {
				return status.Error(codes.InvalidArgument, "attribute key is empty")
			}
			if _reservedAttributesMap[attribute] {
				continue
			}
			if !regexCompiled.MatchString(attribute) {
				return status.Error(codes.InvalidArgument, fmt.Sprintf("event_extensions key %s is invalid.  https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md.  regex used: %s, CloudEvents attribute names MUST consist of lower-case letters ('a' to 'z') or digits ('0' to '9') from the ASCII character set.Attribute names SHOULD be descriptive and terse and SHOULD NOT exceed 20 characters in length.", attribute, pattern))
			}
		}
		keyHint := event.Attributes["partition-key"]
		if keyHint != nil {
			// must be a string
			valueKeyHint := keyHint.GetCeString()
			if fluffycore_utils.IsEmptyOrNil(valueKeyHint) {
				return status.Error(codes.InvalidArgument, "partition-key is empty")
			}
		}
		return nil
	}
	for _, event := range request.Batch.Events {
		if event.Attributes == nil {
			event.Attributes = map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{}
		}
		if fluffycore_utils.IsEmptyOrNil(event) {
			return status.Error(codes.InvalidArgument, "event is empty")
		}
		if fluffycore_utils.IsEmptyOrNil(event.SpecVersion) {
			return status.Error(codes.InvalidArgument, "event.SpecVersion is empty")
		}
		if event.SpecVersion != "1.0" {
			return status.Error(codes.InvalidArgument, "event.SpecVersion is not 1.0")
		}
		if fluffycore_utils.IsEmptyOrNil(event.Type) {
			return status.Error(codes.InvalidArgument, "event.Type is empty")
		}
		if fluffycore_utils.IsEmptyOrNil(event.Source) {
			return status.Error(codes.InvalidArgument, "event.Source is empty")
		}
		if fluffycore_utils.IsEmptyOrNil(event.Id) {
			return status.Error(codes.InvalidArgument, "event.Id is empty")
		}
		err := validateAttributes(event)
		if err != nil {
			return err
		}
		ensureEventTimestamp(event)
	}
	return nil
}

func (s *service) SubmitCloudEvents(ctx context.Context, request *proto_kafkacloudevent.SubmitCloudEventsRequest) (*proto_kafkacloudevent.SubmitCloudEventsResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Interface("request", request).Msg("request")
	err := s.validateCloudEvent(request)
	if err != nil {
		return nil, err
	}
	for _, event := range request.Batch.Events {
		record, err := kafka_franz.CloudEventToKafkaMessage(ctx, event)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert CloudEvent to KafkaMessage")
			return nil, err
		}
		record.Topic = s.config.KafkaConfig.Topic

		var wg sync.WaitGroup
		wg.Add(1)
		keyHint := event.Attributes["partition-key"]
		if keyHint != nil {
			keyHintS := keyHint.GetCeString()
			record.Key = []byte(keyHintS)
		}
		s.KafkaClient.Produce(ctx, record, func(_ *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				log.Error().Err(err).Msg("failed to produce")
			}
		})
		wg.Wait()
	}

	return &proto_kafkacloudevent.SubmitCloudEventsResponse{}, nil
}
