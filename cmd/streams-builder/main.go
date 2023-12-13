package main

import (
	"context"
	"sync"

	"github.com/benthosdev/benthos/v4/public/service"

	// Import only required Benthos components, switch with `components/all` for
	// all standard components.
	_ "github.com/benthosdev/benthos/v4/public/components/io"
	_ "github.com/benthosdev/benthos/v4/public/components/pure"
)

func main() {
	panicOnErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Build the first stream pipeline. Note that we configure each pipeline
	// with its HTTP server disabled as otherwise we would see a port collision
	// when they both attempt to bind to the default address `0.0.0.0:4195`.
	//
	// Alternatively, we could choose to configure each with their own address
	// with the field `http.address`, or we could call `SetHTTPMux` on the
	// builder in order to explicitly override the configured server.
	builderOne := service.NewStreamBuilder()

	err := builderOne.SetYAML(`
http:
  enabled: false

input:
  generate:
    count: 1
    interval: 1ms
    mapping: 'root = "hello world one"'

  processors:
    - mapping: 'root = content().uppercase()'

output:
  stdout: {}
`)
	panicOnErr(err)

	streamOne, err := builderOne.Build()
	panicOnErr(err)

	builderTwo := service.NewStreamBuilder()

	err = builderTwo.SetYAML(`
http:
  enabled: false

input:
  generate:
    
    interval: 1ms
    mapping: 'root = "hello world two"'

  processors:
    - sleep:
        duration: 500ms
    - mapping: 'root = content().capitalize()'

output:
  stdout: {}
`)
	panicOnErr(err)

	streamTwo, err := builderTwo.Build()
	panicOnErr(err)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		panicOnErr(streamOne.Run(context.Background()))
	}()
	go func() {
		defer wg.Done()
		panicOnErr(streamTwo.Run(context.Background()))
	}()
	streamOne.Stop(context.Background())
	wg.Wait()

}
