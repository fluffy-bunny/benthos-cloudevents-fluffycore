package cloudeventoutput

import (
	"context"
	"encoding/json"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	log "github.com/rs/zerolog/log"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

func (s *service) Write(ctx context.Context, message *benthos_service.Message) error {
	content, err := message.AsBytes()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert message to bytes")
		return err
	}
	generic := []map[string]interface{}{}
	err = json.Unmarshal(content, &generic)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal generic")
		return err
	}
	log.Info().Interface("generic", generic).Msg("generic")

	ceBatch := &cloudevents.CloudEventBatch{}
	err = protojson.Unmarshal(content, ceBatch)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CloudEventBatch")
		return err
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
