package runtime

import (
	"context"

	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	contracts_runtime "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/runtime"
	services_bloblang "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/bloblang"
	services_cloudeventoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/cloudeventoutput"
	services_justlogitoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/justlogitoutput"
	services_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/kafkaclient"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	fluffycore_runtime "github.com/fluffy-bunny/fluffycore/runtime"
)

type (
	Startup struct{}
)

func NewStartup() contracts_runtime.IStartup {
	return &Startup{}
}
func (s *Startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	config := &contracts_config.Config{}
	err := fluffycore_runtime.LoadConfig(&fluffycore_contracts_runtime.ConfigOptions{
		Destination: config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
	})
	if err != nil {
		panic(err)
	}
	di.AddInstance[*contracts_config.Config](builder, config)
	services_cloudeventoutput.AddSingletonCloudEventOutput(builder)
	services_justlogitoutput.AddSingletonJustLogItOutput(builder)
	services_kafkaclient.AddSingletonKafkaDeadLetterClient(builder)
	services_bloblang.AddSingletonBlobLangFuncs(builder)
}
