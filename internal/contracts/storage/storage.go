package storage

import (
	centrifuge "github.com/centrifugal/centrifuge-go"
)

type (
	StoreStreamPostitionRequest struct {
		StreamPosition *centrifuge.StreamPosition
	}
	StoreStreamPostitionResponse     struct{}
	GetLatestStreamPostitionRequest  struct{}
	GetLatestStreamPostitionResponse struct {
		StreamPosition *centrifuge.StreamPosition
	}

	// ICentrifugeInputStorage ...
	ICentrifugeInputStorage interface {
		StoreStreamPostition(request *StoreStreamPostitionRequest) (*StoreStreamPostitionResponse, error)
		GetLatestStreamPostition(request *GetLatestStreamPostitionRequest) (*GetLatestStreamPostitionResponse, error)
	}
)
