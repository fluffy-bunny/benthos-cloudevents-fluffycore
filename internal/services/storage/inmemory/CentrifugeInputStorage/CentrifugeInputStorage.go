package CentrifugeInputStorage

import (
	"sync"

	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	status "github.com/gogo/status"
	codes "google.golang.org/grpc/codes"
)

type (
	service struct {
		mutex           sync.Mutex
		streamPositions map[string]*centrifuge.StreamPosition
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_storage.ICentrifugeInputStorage = (*service)(nil)
}

func (s *service) Ctor() contracts_storage.ICentrifugeInputStorage {
	return &service{
		streamPositions: make(map[string]*centrifuge.StreamPosition),
	}
}

func AddSingletonCentrifugeInputStorage(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_storage.ICentrifugeInputStorage](cb, stemService.Ctor)
}
func (s *service) validateStoreStreamPostitionRequest(request *contracts_storage.StoreStreamPostitionRequest) error {
	if fluffycore_utils.IsNil(request) {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Namespace) {
		return status.Error(codes.InvalidArgument, "request.Namespace is empty")
	}
	if fluffycore_utils.IsNil(request.StreamPosition) {
		return status.Error(codes.InvalidArgument, "request.StreamPosition is nil")
	}
	if fluffycore_utils.IsNil(request.StreamPosition.Epoch) {
		return status.Error(codes.InvalidArgument, "request.StreamPosition.Epoch is nil")
	}
	return nil
}
func (s *service) StoreStreamPostition(request *contracts_storage.StoreStreamPostitionRequest) (*contracts_storage.StoreStreamPostitionResponse, error) {
	if err := s.validateStoreStreamPostitionRequest(request); err != nil {
		return nil, err
	}
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
	s.streamPositions[request.Namespace] = request.StreamPosition
	return nil, nil
}
func (s *service) validateGetLatestStreamPostitionRequest(request *contracts_storage.GetLatestStreamPostitionRequest) error {
	if fluffycore_utils.IsNil(request) {
		return status.Error(codes.InvalidArgument, "request is nil")
	}
	if fluffycore_utils.IsEmptyOrNil(request.Namespace) {
		return status.Error(codes.InvalidArgument, "request.Namespace is empty")
	}
	return nil
}
func (s *service) GetLatestStreamPostition(request *contracts_storage.GetLatestStreamPostitionRequest) (*contracts_storage.GetLatestStreamPostitionResponse, error) {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
	if err := s.validateGetLatestStreamPostitionRequest(request); err != nil {
		return nil, err
	}
	// ok to be nil
	streamPosition := s.streamPositions[request.Namespace]
	if streamPosition == nil {
		streamPosition = &centrifuge.StreamPosition{
			Epoch:  "FMvS",
			Offset: 0,
		}

	}
	return &contracts_storage.GetLatestStreamPostitionResponse{
		Namespace:      request.Namespace,
		StreamPosition: streamPosition,
	}, nil

}
