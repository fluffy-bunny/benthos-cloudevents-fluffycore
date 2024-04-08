package storage

import (
	centrifuge "github.com/centrifugal/centrifuge-go"
)

type (
	StoreStreamPostitionRequest struct {
		Namespace      string
		StreamPosition *centrifuge.StreamPosition
	}
	StoreStreamPostitionResponse struct {
	}
	GetLatestStreamPostitionRequest struct {
		Namespace string
	}
	GetLatestStreamPostitionResponse struct {
		Namespace      string
		StreamPosition *centrifuge.StreamPosition
	}

	// ICentrifugeInputStorage ...
	ICentrifugeInputStorage interface {
		StoreStreamPostition(request *StoreStreamPostitionRequest) (*StoreStreamPostitionResponse, error)
		GetLatestStreamPostition(request *GetLatestStreamPostitionRequest) (*GetLatestStreamPostitionResponse, error)
	}
)
