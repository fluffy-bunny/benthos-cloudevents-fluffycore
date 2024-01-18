package bloblang

import (
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_bloblang "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/bloblang"
	contracts_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/kafkaclient"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	service struct {
		deadLetterClient contracts_kafkaclient.IDeadLetterClient
	}
)

var stemService = &service{}

func (s *service) Ctor(deadLetterClient contracts_kafkaclient.IDeadLetterClient) *service {
	return &service{
		deadLetterClient: deadLetterClient,
	}
}
func init() {
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
	var _ contracts_bloblang.IBlobLangFuncs = (*service)(nil)
}
func AddSingletonBlobLangFuncs(cb di.ContainerBuilder) {
	di.AddSingleton[*service](cb, stemService.Ctor,
		contracts_benthos.TypeIBenthosRegistration,
		contracts_bloblang.TypeIBlobLangFuncs)
}
