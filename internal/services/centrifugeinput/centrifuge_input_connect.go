package centrifugeinput

import (
	"context"
	"sync"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

// Establish a connection to the upstream service. Connect will always be
// called first when a reader is instantiated, and will be continuously
// called with back off until a nil error is returned.
//
// The provided context remains open only for the duration of the connecting
// phase, and should not be used to establish the lifetime of the connection
// itself.
//
// Once Connect returns a nil error the Read method will be called until
// either ErrNotConnected is returned, or the reader is closed.
func (s *service) Connect(ctx context.Context) error {
	log := s.log
	// this call is to make sure we can get to our storage.
	_, err := s.centrifugeInputStorage.
		GetLatestStreamPostition(&contracts_storage.GetLatestStreamPostitionRequest{
			Namespace: s.channel,
		})
	if err != nil {
		log.Error().Err(err).Msg("failed to GetLatestStreamPostition")
		return err
	}

	// we may get multiple connect calls since this is benthos doing it.
	// if we already have a subscription don't create another one.
	if s.sub != nil {
		return nil
	}

	client, err := s.centrifugeClient.GetClient()
	if err != nil {
		log.Error().Err(err).Msg("failed to GetClient")
		return err
	}

	sub, err := client.NewSubscription(s.channel,
		centrifuge.SubscriptionConfig{
			Recoverable: true,
			JoinLeave:   true,
			Positioned:  true,
		})
	if err != nil {
		log.Error().Err(err).Msg("failed to NewSubscription")
		return err
	}
	s.sub = sub
	// register for all the events
	sub.OnError(s.OnSubscriptionErrorHandler)
	sub.OnJoin(s.OnJoinHandler)
	sub.OnLeave(s.OnLeaveHandler)
	sub.OnPublication(s.OnPublicationHandler)
	sub.OnSubscribed(s.OnSubscribedHandler)
	sub.OnSubscribing(s.OnSubscribingHandler)
	sub.OnUnsubscribed(s.OnUnsubscribedHandler)

	err = sub.Subscribe()
	if err != nil {
		log.Error().Err(err).Msg("failed to Subscribe")
		return err
	}
	// take a READ lock where our publish will unlock it upon getting a new message
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--
	s.mutexBentosRead.Lock()
	// we block the read until the dispatched message has been acknowledged
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--

	return nil
}

// Read a single message from a source, along with a function to be called
// once the message can be either acked (successfully sent or intentionally
// filtered) or nacked (failed to be processed or dispatched to the output).
//
// The AckFunc will be called for every message at least once, but there are
// no guarantees as to when this will occur. If your input implementation
// doesn't have a specific mechanism for dealing with a nack then you can
// wrap your input implementation with AutoRetryNacks to get automatic
// retries.
//
// If this method returns ErrNotConnected then Read will not be called again
// until Connect has returned a nil error. If ErrEndOfInput is returned then
// Read will no longer be called and the pipeline will gracefully terminate.
func (s *service) Read(ctx context.Context) (*benthos_service.Message, benthos_service.AckFunc, error) {

	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--
	s.mutexBentosRead.Lock()
	// we block the read until the dispatched message has been acknowledged
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--

	// create the ACK Func.  Only then will we commit back to our store where our offset is.
	type ackTracker struct {
		Done                       bool
		mutex                      sync.Mutex
		dispatchedPublicationEvent *PublicationEvent
	}
	ac := &ackTracker{}
	ackFunc := func(ctx context.Context, ackErr error) error {
		// Benthos can call this more than once.
		//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--
		ac.mutex.Lock()
		defer ac.mutex.Unlock()
		//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--
		if ac.Done {
			return nil
		}

		_, err := s.centrifugeInputStorage.StoreStreamPostition(
			&contracts_storage.StoreStreamPostitionRequest{
				Namespace:      s.channel,
				StreamPosition: ac.dispatchedPublicationEvent.streamPosition,
			})
		if err != nil {
			log.Error().Err(err).Msg("failed to StoreStreamPostition")
			return err
		}
		if ackErr != nil {
			log.Error().Err(ackErr).Msg("failed to ack")
			return ackErr
		}
		log.Debug().Msg("ack")
		// ok the message has been acked so we release the semaphore
		// this will let a publish happen again and we then can release the read lock
		s.releaseSemaphore()
		ac.Done = true
		return nil
	}

	content := []byte(s.currentPublicationEvent.publicationEvent.Data)

	msg := benthos_service.NewMessage(content)
	ac.dispatchedPublicationEvent = s.currentPublicationEvent
	s.currentPublicationEvent = nil
	return msg, ackFunc, nil
}
func (s *service) Close(ctx context.Context) error {
	log := zerolog.Ctx(ctx).With().Logger()
	if s.sub != nil {
		err := s.sub.Unsubscribe()
		if err != nil {
			log.Error().Err(err).Msg("failed to Unsubscribe")
			return err
		}
		s.sub = nil
	}
	return s.centrifugeClient.Close(ctx)
}
func (s *service) OnConnectedHandler(e centrifuge.ConnectedEvent) {
	log.Info().Msg("OnConnectedHandler")
}

func (s *service) OnSubscribingHandler(centrifuge.SubscribingEvent) {
	log.Info().Msg("OnSubscribingHandler")
}

func (s *service) OnSubscribedHandler(e centrifuge.SubscribedEvent) {
	log.Info().Msg("OnSubscribedHandler")
	// here we establish where centrifuge is in the stream
	// pull the stream position and store it.
	s.subscribedStreamPosition = e.StreamPosition
	// fire up a go routine to catch up with history
	s.goCatchupHistory()
}
func (s *service) goCatchupHistory() {
	// that a lock no matter what, we will release it when we have caught up with history.
	s.wgPublsh.Add(1)
	// tell our current go routine, if running, to stop
	s.stopHistoryCatchup = true
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// we only allow one go routine to catch up with history at a time
	s.mutexHistoryCatchup.Lock()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~

	// reset our stop flag
	s.stopHistoryCatchup = false
	go func() {
		// once caught up we release the mutex
		defer func() {
			// release our wait group reference
			s.wgPublsh.Done()
			// unlock the history catch up mutex
			s.mutexHistoryCatchup.Unlock()
		}()
		log := log.With().Caller().Interface("subscribedStreamPosition", s.subscribedStreamPosition).Logger()
		getLatestStreamPostitionResponse, err := s.centrifugeInputStorage.
			GetLatestStreamPostition(&contracts_storage.GetLatestStreamPostitionRequest{
				Namespace: s.channel,
			})
		if err != nil {
			log.Error().Err(err).Msg("failed to GetLatestStreamPostition")
			return
		}
		currentStreamPosition := getLatestStreamPostitionResponse.StreamPosition

		log.Debug().Interface("currentStreamPosition", currentStreamPosition).Msg("currentStreamPosition - initial")
		if currentStreamPosition == nil {
			// nothing to but to write the one we got from the subscription
			_, err := s.centrifugeInputStorage.StoreStreamPostition(
				&contracts_storage.StoreStreamPostitionRequest{
					Namespace:      s.channel,
					StreamPosition: s.subscribedStreamPosition,
				})
			if err != nil {
				log.Error().Err(err).Msg("failed to StoreStreamPostition")
			}
			return
		}
		if currentStreamPosition.Offset == s.subscribedStreamPosition.Offset {
			// we are caught up
			return
		}
		currentStreamPosition.Epoch = s.subscribedStreamPosition.Epoch
		// stream position will get updated downstream after we get an ACK
		for {
			if s.stopHistoryCatchup {
				break
			}
			getLatestStreamPostitionResponse, err := s.centrifugeInputStorage.
				GetLatestStreamPostition(&contracts_storage.GetLatestStreamPostitionRequest{
					Namespace: s.channel,
				})
			if err != nil {
				log.Error().Err(err).Msg("failed to GetLatestStreamPostition")
				return
			}
			currentStreamPosition := getLatestStreamPostitionResponse.StreamPosition
			log.Debug().Interface("currentStreamPosition", currentStreamPosition).Msg("currentStreamPosition - loop")

			if currentStreamPosition.Offset == s.subscribedStreamPosition.Offset {
				log.Info().Msg("we are caught up")
				// we are caught up
				return
			}
			historyResult, err := s.sub.History(context.Background(),
				centrifuge.WithHistorySince(currentStreamPosition),
				centrifuge.WithHistoryLimit(100),
			)
			if err != nil {
				log.Printf("error getting history: %v", err)
				break
			}
			if len(historyResult.Publications) == 0 {
				break
			}
			for _, publication := range historyResult.Publications {
				if s.stopHistoryCatchup {
					break
				}
				s.internalOnPublicationHandler(&centrifuge.PublicationEvent{
					Publication: publication,
				})
			}

		}
	}()
}
func (s *service) internalOnPublicationHandler(e *centrifuge.PublicationEvent) {
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// only one message at at time is allowed to come in from centrifuge
	err := s.acquireSemaphore(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to acquireSemaphore")
	}
	// don't release the semaphore here, it has to be released when the benthos consumer ACKS the message
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	s.currentPublicationEvent = &PublicationEvent{
		publicationEvent: e,
		streamPosition: &centrifuge.StreamPosition{
			Epoch:  s.subscribedStreamPosition.Epoch,
			Offset: e.Offset,
		},
	}
	// we now can let the bentos read happen
	s.mutexBentosRead.Unlock()
}

func (s *service) OnPublicationHandler(e centrifuge.PublicationEvent) {
	log.Debug().Msg("OnPublicationHandler")
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	// we are blocking publication to use from centrifuge until we have caught up with history
	s.wgPublsh.Wait()
	//--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~
	s.internalOnPublicationHandler(&e)
}
