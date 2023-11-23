package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	// Import only pure and standard io Benthos components
	_ "github.com/benthosdev/benthos/v4/public/components/io"
	_ "github.com/benthosdev/benthos/v4/public/components/pure"

	_ "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/bloblang"

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

	startup := internal_runtime.NewStartup()
	builder := di.Builder()
	startup.ConfigureServices(ctx, builder)
	internal.Container = builder.Build()

	registrations := di.Get[[]contracts_benthos.IBenthosRegistration](internal.Container)
	for _, registration := range registrations {
		registration.Register()
	}
	waitChannel := make(chan os.Signal, 1)

	originalArgs := os.Args

	// cancel context
	ctx, cancel := context.WithCancel(ctx)
	// benthos thinks its the only one, so lets replace the args and then set them back when it launches.
	benthosOSArgsS := os.Getenv("BENTHOS_OS_ARGS")
	fmt.Println("BENTHOS_OS_ARGS", benthosOSArgsS)
	// split them
	benthosOSArgs := strings.Split(benthosOSArgsS, ",")
	os.Args = []string{
		originalArgs[0],
	}
	os.Args = append(os.Args, benthosOSArgs...)
	log.Info().Interface("os.Args", os.Args).Msg("os.Args")
	os.Args[0] = "benthos"
	go func() {
		service.RunCLI(context.Background())
	}()
	time.Sleep(5 * time.Second)
	os.Args = originalArgs
	// do my stuff
	signal.Notify(
		waitChannel,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	<-waitChannel
	cancel()
}
