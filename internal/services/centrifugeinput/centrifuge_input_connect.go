package centrifugeinput

import (
	"context"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
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
