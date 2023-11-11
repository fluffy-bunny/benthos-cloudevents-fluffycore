package runtime

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	IStartup interface {
		ConfigureServices(ctx context.Context, builder di.ContainerBuilder)
	}
)
