package runtime

import (
	"context"

	contracts_runtime "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/runtime"
	services_cloudeventoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/services/cloudeventoutput"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	Startup struct{}
)

func NewStartup() contracts_runtime.IStartup {
	return &Startup{}
}
func (s *Startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	services_cloudeventoutput.AddSingletonCloudEventOutput(builder)
}
