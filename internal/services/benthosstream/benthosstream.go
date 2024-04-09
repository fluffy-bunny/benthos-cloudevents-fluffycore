package benthosstream

import (
	"context"
	"os"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	"github.com/rs/zerolog"
)

type (
	service struct {
		stream *benthos_service.Stream
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_benthos.IBenthosStream = (*service)(nil)
}

func AddTransientBenthosStream(builder di.ContainerBuilder) {
	di.AddTransient[contracts_benthos.IBenthosStream](builder, stemService.Ctor)
}

func (s *service) Ctor() (contracts_benthos.IBenthosStream, error) {
	return &service{}, nil
}

func (s *service) Configure(ctx context.Context, config *contracts_benthos.CentrifugeStreamConfig) (err error) {
	log := zerolog.Ctx(ctx).With().Logger()
	builderOne := benthos_service.NewStreamBuilder()
	// load yaml from file
	content, err := os.ReadFile(config.ConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read file")
		return err
	}
	err = builderOne.SetYAML(string(content))
	if err != nil {
		log.Error().Err(err).Msg("failed to set yaml")
		return err
	}
	stream, err := builderOne.Build()
	if err != nil {
		return err
	}
	s.stream = stream
	return nil
}

func (s *service) Run(ctx context.Context) (err error) {
	return s.stream.Run(ctx)
}
func (s *service) Stop(ctx context.Context) (err error) {
	return s.stream.Stop(ctx)
}
