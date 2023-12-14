package bloblang

import (
	"github.com/benthosdev/benthos/v4/public/bloblang"
)

func (s *service) Register() error {
	deadLetterItSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewAnyParam("headers")).
		Param(bloblang.NewAnyParam("content"))
	err := bloblang.RegisterFunctionV2("deadletterit", deadLetterItSpec, s.DeadLetterIt)
	if err != nil {
		return err
	}
	wrapItSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewAnyParam("headers")).
		Param(bloblang.NewAnyParam("content"))
	err = bloblang.RegisterFunctionV2("wrapit", wrapItSpec, s.WrapIt)
	if err != nil {
		return err
	}
	errorLogItSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewAnyParam("headers")).
		Param(bloblang.NewAnyParam("content"))
	err = bloblang.RegisterFunctionV2("errorlogit", errorLogItSpec, s.ErrorLogIt)
	if err != nil {
		return err
	}

	return err
}
