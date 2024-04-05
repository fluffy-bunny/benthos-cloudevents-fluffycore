package runtime

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	contracts_runtime "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/runtime"
	services_bloblang "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/bloblang"
	services_centrifugeclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/centrifugeclient"
	services_centrifugeinput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/centrifugeinput"
	services_centrifugestream "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/centrifugestream"
	services_cloudeventoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/cloudeventoutput"
	services_justlogitoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/justlogitoutput"
	services_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/kafkaclient"
	services_kafkasaslstream "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/kafkasaslstream"
	services_kafkastream "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/kafkastream"
	services_storage_inmemory_CentrifugeInputStorage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/storage/inmemory/CentrifugeInputStorage"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	fluffycore_runtime "github.com/fluffy-bunny/fluffycore/runtime"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

type (
	Startup struct{}
)

func NewStartup() contracts_runtime.IStartup {
	return &Startup{}
}
func (s *Startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	log := log.With().Caller().Str("method", "GetConfigOptions").Logger()

	// here we load a config file and merge it over the default.
	initialConfigOptions := &fluffycore_contracts_runtime.ConfigOptions{
		Destination: &contracts_config.InitialConfig{},
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	err := fluffycore_runtime.LoadConfig(initialConfigOptions)
	if err != nil {
		panic(err)
	}

	err = onLoadFileConfig(context.Background(), initialConfigOptions.Destination.(*contracts_config.InitialConfig).ConfigFiles.RootPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to onLoadFileConfig")
		panic(err)
	}

	defaultConfig := &contracts_config.Config{}
	err = json.Unmarshal([]byte(contracts_config.ConfigDefaultJSON), defaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
	}
	log.Info().Interface("defaultConfig", defaultConfig).Msg("config after merge")

	config := &contracts_config.Config{}
	err = fluffycore_runtime.LoadConfig(&fluffycore_contracts_runtime.ConfigOptions{
		Destination: config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
	})
	if err != nil {
		panic(err)
	}
	log.Info().Interface("final_config", config).Msg("config after merge")

	di.AddInstance[*contracts_config.Config](builder, config)
	services_centrifugeinput.AddSingletonCentrifugeInput(builder)
	services_centrifugeclient.AddSingletonCentrifugeClient(builder)
	services_cloudeventoutput.AddSingletonCloudEventOutput(builder)
	services_justlogitoutput.AddSingletonJustLogItOutput(builder)
	services_kafkaclient.AddSingletonKafkaDeadLetterClient(builder)
	services_bloblang.AddSingletonBlobLangFuncs(builder)
	if false {
		if config.EnableKafkaSASL {
			services_kafkasaslstream.AddSingletonIBenthosStream(builder)
		} else {
			services_kafkastream.AddSingletonIBenthosStream(builder)
		}

	}
	services_centrifugestream.AddSingletonIBenthosStream(builder)
	services_storage_inmemory_CentrifugeInputStorage.AddSingletonCentrifugeInputStorage(builder)
}

// onLoadMastodonConfig will load a file and merge it over the default config
// you can still use ENV variables to replace as well.  i.e. for secrets that only come in that way.
// ---------------------------------------------------------------------
func onLoadFileConfig(ctx context.Context, configFilePath string) error {
	log := zerolog.Ctx(ctx).With().Str("method", "OnConfigureServicesLoadIDPs").Logger()
	log.Info().Str("mastodonConfigFilePath", configFilePath).Msg("loading mastodon config")

	if fluffycore_utils.IsEmptyOrNil(configFilePath) {
		// may not be an error
		log.Info().Msg("onLoadFileConfig did not load a file due to empty configFilePath")
		return nil
	}

	fileContent, err := os.ReadFile(configFilePath)

	if err != nil {
		log.Error().Err(err).Msg("failed to read mastodonConfigFilePath")
		return err
	}
	fixedFileContent := fluffycore_utils.ReplaceEnv(string(fileContent), "${%s}")
	overlay := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(fixedFileContent)).Decode(&overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ragePath")
		return err
	}
	src := map[string]interface{}{}

	err = json.NewDecoder(strings.NewReader(string(contracts_config.ConfigDefaultJSON))).Decode(&src)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
		return err
	}
	err = fluffycore_utils.ReplaceMergeMap(overlay, src)
	if err != nil {
		log.Error().Err(err).Msg("failed to ReplaceMergeMap")
		return err
	}
	bb, err := json.Marshal(overlay)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal overlay")
		return err
	}
	contracts_config.ConfigDefaultJSON = bb

	return nil

}
