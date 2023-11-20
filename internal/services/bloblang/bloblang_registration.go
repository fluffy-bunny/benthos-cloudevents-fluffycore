package bloblang

import (
	"github.com/benthosdev/benthos/v4/public/bloblang"
)

func (s *service) Register() error {
	crazyObjectSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewAnyParam("headers")).
		Param(bloblang.NewAnyParam("content"))
	err := bloblang.RegisterFunctionV2("deadletterit", crazyObjectSpec, s.DeadLetterIt)

	return err
}
