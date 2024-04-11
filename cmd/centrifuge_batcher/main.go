package main

import (
	"context"
	"os"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	pkg_contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/contracts/centrifuge"
	pkg_services_centrifuge_batcher "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/services/centrifuge/batcher"
	pkg_services_centrifuge_client "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/services/centrifuge/client"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	ctx := log.WithContext(context.Background())
	var cancelCtx context.CancelFunc

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	ctx, cancelCtx = context.WithCancel(ctx)

	builder := di.Builder()
	pkg_services_centrifuge_batcher.AddTransientCentrifugeStreamBatcher(builder)
	pkg_services_centrifuge_client.AddTransientCentrifugeClient(builder)
	container := builder.Build()

	centrifugeStreamBatcher := di.Get[pkg_contracts_centrifuge.ICentrifugeStreamBatcher](container)
	centrifugeStreamBatcher.Configure(ctx, &pkg_contracts_centrifuge.CentrifugeConfig{
		Channel:         "chat:index",
		BatchSize:       3,
		NumberOfBatches: 2,
		HistoricalStreamPosition: &centrifuge.StreamPosition{
			Epoch:  "dddd",
			Offset: 0,
		},
		CentrifugeClientConfig: &pkg_contracts_centrifuge.CentrifugeClientConfig{
			Endpoint: "ws://localhost:8079/connection/websocket",
			GetToken: func(e centrifuge.ConnectionTokenEvent) (string, error) {
				return pkg_services_centrifuge_client.ExampleConnToken("49", 0), nil
			},
		},
	})

	err := centrifugeStreamBatcher.Start(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("centrifugeStreamBatcher.Start")
	}

	count := 0
	for {
		if count > 10000 {
			cancelCtx()
			break
		}
		time.Sleep(5 * time.Second)

		batch, err := centrifugeStreamBatcher.GetBatch(ctx, true)
		if err != nil {
			log.Error().Err(err).Msg("centrifugeStreamBatcher.GetBatch")
		}
		log.Info().Interface("batch", batch).Msg("centrifugeStreamBatcher.GetBatch")

		count++
	}
	count = 0
	for {
		if count > 10 {
			break
		}
		time.Sleep(30 * time.Second)
	}
	err = centrifugeStreamBatcher.Stop(ctx)
	if err != nil {
		log.Error().Err(err).Msg("centrifugeStreamBatcher.Stop")
	}
}
