package main

import (
	internal_runtime "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/runtime"
	internal_version "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/version"
	fluffycore_cobracore_cmd "github.com/fluffy-bunny/fluffycore/cobracore/cmd"
)

func main() {
	startup := internal_runtime.NewStartup()
	fluffycore_cobracore_cmd.SetVersion(internal_version.Version())
	fluffycore_cobracore_cmd.Execute(startup)
}
