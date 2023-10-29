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
	_ "github.com/benthosdev/benthos/v4/public/components/all"

	_ "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/input"
	_ "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/output"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// create a logger and add it to the context
	logz := zerolog.New(os.Stdout).With().Timestamp().Logger()

	log.Logger = logz.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("main")
	service.RunCLI(context.Background())
}
