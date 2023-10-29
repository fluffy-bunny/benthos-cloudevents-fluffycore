package runtime

import (
	"context"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_middleware "github.com/fluffy-bunny/fluffycore/contracts/middleware"
	zerolog "github.com/rs/zerolog"

	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
)

type (
	startup struct {
		fluffycore_contracts_runtime.UnimplementedStartup
		RootContainer di.Container
		config        *contracts_config.Config
		configOptions *fluffycore_contracts_runtime.ConfigOptions
	}
)

func NewStartup() fluffycore_contracts_runtime.IStartup {
	return &startup{}
}
func (s *startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	log.Info().Msg("ConfigureServices")
}
func (s *startup) Configure(ctx context.Context, rootContainer di.Container, unaryServerInterceptorBuilder fluffycore_contracts_middleware.IUnaryServerInterceptorBuilder, streamServerInterceptorBuilder fluffycore_contracts_middleware.IStreamServerInterceptorBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()
	log.Info().Msg("Configure")
}
func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	s.config = &contracts_config.Config{}
	s.configOptions = &fluffycore_contracts_runtime.ConfigOptions{
		Destination: s.config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	return s.configOptions
}

func (s *startup) SetRootContainer(container di.Container) {
	s.RootContainer = container
}
