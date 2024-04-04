package centrifugeinput

import (
	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	service struct {
		// grpcUrl: i.e. grpc://localhost:5001
		endpoint string
		// channel: is a hint to the processor.  This allows a processor to have a single app that takes all the requests.
		channel                string
		logger                 *benthos_service.Logger
		centrifugeInputStorage contracts_storage.ICentrifugeInputStorage
	}
)

var stemService = &service{}

func (s *service) Ctor(centrifugeInputStorage contracts_storage.ICentrifugeInputStorage) *service {
	return &service{
		centrifugeInputStorage: centrifugeInputStorage,
	}
}
func init() {
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
}
func AddSingletonCloudEventOutput(cb di.ContainerBuilder) {
	di.AddSingleton[*service](cb, stemService.Ctor,
		contracts_benthos.TypeIBenthosRegistration,
	)
}
