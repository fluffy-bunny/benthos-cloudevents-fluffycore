package centrifugeinput

import (
	"context"
	"sync"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/centrifuge"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

type (
	service struct {
		contracts_centrifuge.UnimplementedSubscriptionHandlers

		// channel: is a hint to the processor.  This allows a processor to have a single app that takes all the requests.
		channel                string
		logger                 *benthos_service.Logger
		centrifugeInputStorage contracts_storage.ICentrifugeInputStorage
		centrifugeClient       contracts_centrifuge.ICentrifugeClient
		streamPosition         *centrifuge.StreamPosition
		log                    zerolog.Logger
		// this semaphore is to block the publish until our downstream has consumed and acked the message
		sem *semaphore.Weighted
		// this mutexPublish is to block the publish until we have caught up with history
		mutexPublish sync.Mutex

		// this  mutexRead is to guard against multiple reads
		mutexRead sync.Mutex

		currentPublicationEvent *centrifuge.PublicationEvent

		sub *centrifuge.Subscription
	}
)

var stemService = &service{}

func (s *service) Ctor(
	centrifugeInputStorage contracts_storage.ICentrifugeInputStorage,
	centrifugeClient contracts_centrifuge.ICentrifugeClient) *service {

	log := log.With().Caller().Str("input", InputName).Logger()
	sem := semaphore.NewWeighted(1)
	svc := &service{
		centrifugeInputStorage: centrifugeInputStorage,
		centrifugeClient:       centrifugeClient,
		log:                    log,
		sem:                    sem,
	}

	return svc
}
func init() {
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
}
func AddSingletonCentrifugeInput(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_benthos.IBenthosRegistration](cb, stemService.Ctor)
}

func (s *service) acquireSemaphore(ctx context.Context) error {
	err := s.sem.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) releaseSemaphore() {
	s.sem.Release(1)
}
