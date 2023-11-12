package cloudeventoutput

import (
	"context"
	"encoding/json"
	"time"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	benthos_message "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/benthos/message"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	proto_cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"

	status "github.com/gogo/status"
	log "github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type kafkaPayload struct {
	Headers map[string]interface{} `json:"headers"`
	Value   map[string]interface{} `json:"value"`
}

func (s *service) ItemToCloudEvent2(ctx context.Context, bytes *[]byte) (*proto_cloudevents.CloudEvent, error) {
	ce := &proto_cloudevents.CloudEvent{}
	err := protojson.Unmarshal(*bytes, ce)
	if err != nil {
		return nil, err
	}
	return ce, nil
}

func (s *service) ItemToCloudEvent(ctx context.Context, item map[string]interface{}) (*proto_cloudevents.CloudEvent, error) {
	headers, ok := item["headers"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to find headers")
	}
	// is headers a map
	headersMap, ok := headers.(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to cast headers to map[string]interface{}")
	}
	value, ok := item["value"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to find value")
	}
	// is value a map
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to cast value to map[string]interface{}")
	}
	ce := &proto_cloudevents.CloudEvent{
		Attributes: make(map[string]*proto_cloudevents.CloudEvent_CloudEventAttributeValue),
	}
	pullStringFromHeader := func(key string) (string, error) {
		value, ok := headersMap[key]
		if !ok {
			return "", status.Error(codes.InvalidArgument, "failed to find "+key)
		}
		// is value a string
		valueString, ok := value.(string)
		if !ok {
			return "", status.Error(codes.InvalidArgument, "failed to cast "+key+" to string")
		}
		// delete it from the map
		delete(headersMap, key)
		return valueString, nil
	}
	// pull spec from headers
	spec, err := pullStringFromHeader("ce_specversion")
	if err != nil {
		return nil, err
	}
	ce.SpecVersion = spec
	// pull type from headers
	ce.Type, err = pullStringFromHeader("ce_type")
	if err != nil {
		return nil, err
	}
	// pull source from headers
	ce.Source, err = pullStringFromHeader("ce_source")
	if err != nil {
		return nil, err
	}
	// pull id from headers
	ce.Id, err = pullStringFromHeader("ce_id")
	if err != nil {
		return nil, err
	}
	// "ce_time": "2023-11-01T15:52:57.9309064Z",
	ceTime, err := pullStringFromHeader("ce_time")
	if err != nil {
		return nil, err
	}
	// convert string to time
	tt, err := time.Parse(time.RFC3339, ceTime)
	if err != nil {
		return nil, err
	}
	ce.Attributes["time"] = &proto_cloudevents.CloudEvent_CloudEventAttributeValue{
		Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeTimestamp{
			CeTimestamp: timestamppb.New(tt),
		},
	}
	contentType, err := pullStringFromHeader("content-type")
	if err != nil {
		return nil, err
	}
	if contentType != "application/json" {
		return nil, status.Error(codes.InvalidArgument, "content-type must be application/json")
	}
	ce.Attributes["content-type"] = &proto_cloudevents.CloudEvent_CloudEventAttributeValue{
		Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
			CeString: contentType,
		},
	}
	// remove any headers that start with kafka_
	var toBeDeleted []string
	for key := range headersMap {
		if len(key) > 6 && key[0:6] == "kafka_" {
			toBeDeleted = append(toBeDeleted, key)
		}
	}
	for _, key := range toBeDeleted {
		delete(headersMap, key)
	}
	// turn the rest into attributes
	for key, value := range headersMap {
		// is value a string
		valueString, ok := value.(string)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "failed to cast "+key+" to string")
		}
		ce.Attributes[key] = &proto_cloudevents.CloudEvent_CloudEventAttributeValue{
			Attr: &proto_cloudevents.CloudEvent_CloudEventAttributeValue_CeString{
				CeString: valueString,
			},
		}
	}
	// convert value to json
	content, err := json.Marshal(valueMap)
	if err != nil {
		return nil, err
	}
	ce.Data = &proto_cloudevents.CloudEvent_TextData{
		TextData: string(content),
	}
	return ce, nil
}

func (s *service) Write(ctx context.Context, message *benthos_service.Message) error {
	content, err := message.AsBytes()
	if err != nil {
		log.Warn().Err(err).Msg("failed to convert message to AsBytes, don't want it anyway if that happens")
		return nil
	}
	if len(content) == 0 {
		log.Warn().Msg("empty message, don't want it anyway if that happens")
		return nil
	}
	parts, err := benthos_message.DeserializeBytes(content)
	if err != nil {
		log.Warn().Err(err).Msg("failed to DeserializeBytes, don't want it anyway if that happens")
		return nil
	}
	processRequest := &proto_cloudeventprocessor.ProcessCloudEventsRequest{}

	ceBatch := &proto_cloudevents.CloudEventBatch{}
	badParts := [][]byte{}
	for _, part := range parts {
		// must be a json object
		ce, err := s.ItemToCloudEvent2(ctx, &part)
		if err != nil {
			badParts = append(badParts, part)
			continue
		}
		ceBatch.Events = append(ceBatch.Events, ce)

	}
	if len(ceBatch.Events) > 0 {
		processRequest.Batch = ceBatch
	}
	if len(badParts) > 0 {
		processRequest.BadBatch = &proto_cloudeventprocessor.ByteBatch{
			Messages: badParts,
		}
	}
	log.Info().Interface("ceBatch", ceBatch).Msg("ceBatch")
	_, err = s.cloudEventProcessorClient.ProcessCloudEvents(ctx, &proto_cloudeventprocessor.ProcessCloudEventsRequest{
		Batch: ceBatch,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to ProcessCloudEvents")
		return err
	}
	return nil
}
