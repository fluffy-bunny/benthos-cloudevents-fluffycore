package centrifugeinput

import (
	"context"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
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
	getLatestStreamPostitionRequest, err := s.centrifugeInputStorage.GetLatestStreamPostition(&contracts_storage.GetLatestStreamPostitionRequest{})
	if err != nil {
		log.Error().Err(err).Msg("failed to GetLatestStreamPostition")
		return err
	}
	// can be nil
	s.streamPosition = getLatestStreamPostitionRequest.StreamPosition

	client := s.centrifugeClient.GetClient()
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
	// register for all the events
	sub.OnError(s.OnSubscriptionErrorHandler)
	sub.OnJoin(s.OnJoinHandler)
	sub.OnLeave(s.OnLeaveHandler)
	sub.OnPublication(s.OnPublicationHandler)
	sub.OnSubscribed(s.OnSubscribedHandler)
	sub.OnSubscribing(s.OnSubscribingHandler)
	sub.OnUnsubscribed(s.OnUnsubscribedHandler)

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
	return nil, nil, nil
}
func (s *service) Close(ctx context.Context) error {
	return nil
}
func (s *service) OnConnectedHandler(e centrifuge.ConnectedEvent) {
	log.Info().Msg("OnConnectedHandler")
}
