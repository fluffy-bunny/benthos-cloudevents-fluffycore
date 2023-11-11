package cloudevents

import (
	"testing"

	proto_cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	require "github.com/stretchr/testify/require"
	suite "github.com/stretchr/testify/suite"
)

type CloudEventCase struct {
	CloudEvent proto_cloudevents.CloudEvent
	Pass       bool
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type CloudEventsTestSuite struct {
	suite.Suite
	CloudEvents []*CloudEventCase
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CloudEventsTestSuite) SetupTest() {
	suite.CloudEvents = []*CloudEventCase{
		{
			Pass: true,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"partition-key": {
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
					"customcontext": {
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		},
		{
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				// no id
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id: "123",
				// no source
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:     "123",
				Source: "source",
				// no spec version
				Type: "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.3333", // wrong spec version
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				// no type
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"partition-key": {
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeBoolean{
							CeBoolean: true, // wrong type
						},
					},
					"customcontext": {
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"bad-context": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"BadContext": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"thenameiswaytolongitmustbeonly20characters": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"a b": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					"": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		}, {
			Pass: false,
			CloudEvent: proto_cloudevents.CloudEvent{
				Id:          "123",
				Source:      "source",
				SpecVersion: "1.0",
				Type:        "test",
				Data: &proto_cloudevents.CloudEvent_TextData{
					TextData: `{"a":1, "b":"two"}`,
				},
				Attributes: map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue{
					" ": {
						// bad context name
						Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
							CeString: "123",
						},
					},
				},
			},
		},
	}
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *CloudEventsTestSuite) TestExample() {

	for _, ce := range suite.CloudEvents {
		err := ValidateCloudEvent(&ce.CloudEvent)
		if ce.Pass {
			require.NoError(suite.T(), err)
		} else {
			require.Error(suite.T(), err)
		}
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CloudEventsTestSuite))
}
