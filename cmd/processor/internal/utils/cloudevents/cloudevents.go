package cloudevents

import (
	"fmt"

	"regexp"

	proto_cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

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

func ValidateCloudEvent(event *proto_cloudevents.CloudEvent) error {
	ceTimestamp := timestamppb.Now()

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

	return nil
}
