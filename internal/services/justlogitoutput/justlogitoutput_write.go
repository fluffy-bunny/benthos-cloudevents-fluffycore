package justlogitoutput

import (
	"context"
	"encoding/json"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	zerolog_log "github.com/rs/zerolog/log"
)

func (s *service) Write(ctx context.Context, message *benthos_service.Message) error {
	log := zerolog_log.With().Caller().Logger()
	//ctx = log.WithContext(ctx)

	content, err := message.AsBytes()
	if err != nil {
		zerolog_log.Warn().Err(err).Msg("failed to convert message to AsBytes, don't want it anyway if that happens")
		return nil
	}
	if len(content) == 0 {
		zerolog_log.Warn().Msg("empty message, don't want it anyway if that happens")
		return nil
	}
	generic := map[string]interface{}{}
	err = json.Unmarshal(content, &generic)
	if err != nil {
		log.Error().Err(err).Msg("failed to Unmarshal, don't want it anyway if that happens")
		return nil
	}

	log.Info().Interface("generic", generic).Msg("generic")
	return nil
}
