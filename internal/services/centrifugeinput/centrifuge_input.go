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
	semaphore "golang.org/x/sync/semaphore"
)

type (
	PublicationEvent struct {
		publicationEvent *centrifuge.PublicationEvent
		streamPosition   *centrifuge.StreamPosition
	}
	service struct {
		contracts_centrifuge.UnimplementedSubscriptionHandlers

		// channel: is a hint to the processor.  This allows a processor to have a single app that takes all the requests.
		channel                  string
		endpoint                 string
		logger                   *benthos_service.Logger
		centrifugeInputStorage   contracts_storage.ICentrifugeInputStorage
		centrifugeClient         contracts_centrifuge.ICentrifugeClient
		subscribedStreamPosition *centrifuge.StreamPosition
		log                      zerolog.Logger
		// this semaphore is to block the publish until our downstream has consumed and acked the message
		sem *semaphore.Weighted

		// this  mutexBentosRead is to guard against multiple reads
		mutexBentosRead sync.Mutex

		currentPublicationEvent *PublicationEvent

		sub *centrifuge.Subscription
		// this is wait group that will block publishing to happen until we have caught up with out history
		wgPublsh            sync.WaitGroup
		mutexHistoryCatchup sync.Mutex
		stopHistoryCatchup  bool
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
func AddTransientCentrifugeInput(cb di.ContainerBuilder) {
	di.AddTransient[contracts_benthos.IBenthosRegistration](cb, stemService.Ctor)
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
