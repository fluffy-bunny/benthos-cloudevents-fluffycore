package main

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	//	benthos_service_servicetest "github.com/benthosdev/benthos/v4/public/service/servicetest"
	//	benthos_service "github.com/benthosdev/benthos/v4/public/service"

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
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/centrifuge"
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
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

	_ = di.Get[*contracts_config.Config](internal.Container)
	_ = di.Get[contracts_storage.ICentrifugeInputStorage](internal.Container)
	_ = di.Get[contracts_centrifuge.ICentrifugeClient](internal.Container)

	registrations := di.Get[[]contracts_benthos.IBenthosRegistration](internal.Container)
	for _, registration := range registrations {
		registration.Register()
	}

	waitChannel := make(chan os.Signal, 1)

	// cancel context
	ctx, cancel := context.WithCancel(ctx)

	benthosStreams := di.Get[[]contracts_benthos.IBenthosStream](internal.Container)
	for _, benthosStream := range benthosStreams {
		go func(benthosStream contracts_benthos.IBenthosStream) {
			err := benthosStream.Run(ctx)
			if err != nil {
				log.Fatal().Err(err).Msg("Error starting stream")
			}
		}(benthosStream)
	}

	// do my stuff
	signal.Notify(
		waitChannel,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	<-waitChannel
	for _, benthosStream := range benthosStreams {
		benthosStream.Stop(ctx)
	}

	cancel()
}

func LoadYamlFile(filename string) (string, error) {
	// load kafka.yml into a string
	kafkaYaml, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(kafkaYaml), nil
}
func FixupFromEnv(str string) string {
	re := regexp.MustCompile(`\${.*?}`)
	matches := re.FindAllString(str, -1)
	for _, match := range matches {
		envVar := os.Getenv(match[2 : len(match)-1])
		str = regexp.MustCompile(regexp.QuoteMeta(match)).ReplaceAllString(str, envVar)
	}
	return str
}
