package centrifugeinput

import (
	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/centrifuge"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

type (
	service struct {
		contracts_centrifuge.UnimplementedSubscriptionHandlers
		// grpcUrl: i.e. grpc://localhost:5001
		endpoint string
		// channel: is a hint to the processor.  This allows a processor to have a single app that takes all the requests.
		channel                string
		logger                 *benthos_service.Logger
		centrifugeInputStorage contracts_storage.ICentrifugeInputStorage
		centrifugeClient       contracts_centrifuge.ICentrifugeClient
		streamPosition         *centrifuge.StreamPosition
		log                    zerolog.Logger
	}
)

var stemService = &service{}

func (s *service) Ctor(
	centrifugeInputStorage contracts_storage.ICentrifugeInputStorage,
	centrifugeClient contracts_centrifuge.ICentrifugeClient) *service {

	log := log.With().Caller().Str("input", InputName).Logger()
	return &service{
		centrifugeInputStorage: centrifugeInputStorage,
		centrifugeClient:       centrifugeClient,
		log:                    log,
	}
}
func init() {
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
}
func AddSingletonCentrifugeInput(cb di.ContainerBuilder) {
	di.AddSingleton[*service](cb, stemService.Ctor,
		contracts_benthos.TypeIBenthosRegistration,
	)
}
