package justlogitoutput

import (
	"context"
	"encoding/json"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	benthos_message "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/benthos/message"
	zerolog_log "github.com/rs/zerolog/log"
)

func (s *service) Write(ctx context.Context, message *benthos_service.Message) error {
	log := zerolog_log.With().Caller().Logger()
	ctx = log.WithContext(ctx)

	content, err := message.AsBytes()
	if err != nil {
		zerolog_log.Warn().Err(err).Msg("failed to convert message to AsBytes, don't want it anyway if that happens")
		return nil
	}
	if len(content) == 0 {
		zerolog_log.Warn().Msg("empty message, don't want it anyway if that happens")
		return nil
	}
	parts, err := benthos_message.DeserializeBytes(content)
	if err != nil {
		zerolog_log.Warn().Err(err).Msg("failed to DeserializeBytes, don't want it anyway if that happens")
		return nil
	}
	genericParts := []map[string]interface{}{}
	for _, part := range parts {
		generic := make(map[string]interface{})
		err = json.Unmarshal(part, &generic)
		if err != nil {
			log.Error().Err(err).Msg("failed to Unmarshal, don't want it anyway if that happens")
			return err
		}
		genericParts = append(genericParts, generic)
	}

	log.Info().Interface("genericParts", genericParts).Msg("generic")
	return nil
}
