package justlogitoutputbatch

import (
	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_justlogitoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/justlogitoutput"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	service struct {
		contracts_benthos.UnimplementedBenthosOutput
		logger *benthos_service.Logger
	}
)

var stemService = &service{}

func (s *service) Ctor() *service {
	return &service{}
}
func init() {
	var _ contracts_justlogitoutput.IJustLogItOutput = (*service)(nil)
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
}
func AddSingletonJustLogItOutput(cb di.ContainerBuilder) {
	di.AddSingleton[*service](cb, stemService.Ctor,
		contracts_benthos.TypeIBenthosRegistration,
		contracts_justlogitoutput.TypeIJustLogItOutput)
}
