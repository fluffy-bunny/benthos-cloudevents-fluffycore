package main

import (
	"context"
	"os"

	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	// Import only pure and standard io Benthos components
	_ "github.com/benthosdev/benthos/v4/public/components/io"
	_ "github.com/benthosdev/benthos/v4/public/components/pure"

	// In order to import _all_ Benthos components for third party services
	// uncomment the following line:
	//_ "github.com/benthosdev/benthos/v4/public/components/all"
	_ "github.com/benthosdev/benthos/v4/public/components/kafka"
	internal "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	_ "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/input"
	_ "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/output"
	internal_runtime "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/runtime"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

func main() {
	ctx := context.Background()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// create a logger and add it to the context
	logz := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	log.Logger = logz.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx = log.Logger.WithContext(ctx)

	log.Info().Msg("main")
	startup := internal_runtime.NewStartup()
	builder := di.Builder()
	startup.ConfigureServices(ctx, builder)
	internal.Container = builder.Build()

	registrations := di.Get[[]contracts_benthos.IBenthosRegistration](internal.Container)
	for _, registration := range registrations {
		registration.Register()
	}
	service.RunCLI(context.Background())
}
