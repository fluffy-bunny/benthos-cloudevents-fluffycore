package centrifuge

import (
	"context"

	centrifuge "github.com/centrifugal/centrifuge-go"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
)

type (
	OnBatchReady           func()
	CentrifugeClientConfig struct {
		Endpoint string `json:"endpoint"`
		GetToken func(evt centrifuge.ConnectionTokenEvent) (string, error)
	}
	CentrifugeConfig struct {
		Channel                  string                     `json:"channel"`
		BatchSize                int32                      `json:"batchSize"`
		NumberOfBatches          int                        `json:"numberOfBatches"`
		HistoricalStreamPosition *centrifuge.StreamPosition `json:"historicalStreamPosition"`
		OnBatchReady             OnBatchReady               `json:"-"`
		CentrifugeClientConfig   *CentrifugeClientConfig    `json:"centrifugeClientConfig"`
	}
	Batch struct {
		Publications        []centrifuge.Publication
		LastStreamPostition *centrifuge.StreamPosition
	}
	BatchState int

	// ICentrifugeTokenAccessor is exected to be a singleton
	ISingletonCentrifugeTokenAccessor interface {
		GetToken(evt centrifuge.ConnectionTokenEvent) (string, error)
	}
	ICentrifugeStreamBatcher interface {
		Configure(ctx context.Context, config *CentrifugeConfig) error
		GetBatch(ctx context.Context, flush bool) (*Batch, error)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		IsRunning() bool
		BatchState() BatchState
	}
	ICentrifugeClient interface {
		fluffycore_contracts_common.IDispose
		Configure(ctx context.Context, config *CentrifugeClientConfig)
		GetClient() (*centrifuge.Client, error)
	}
	IClientHandlers interface {
		centrifuge.ConnectedHandler
		centrifuge.ConnectingHandler
		centrifuge.DisconnectHandler
		centrifuge.MessageHandler
		centrifuge.ServerPublicationHandler
		centrifuge.ServerSubscribedHandler
		centrifuge.ServerSubscribingHandler
		centrifuge.ServerUnsubscribedHandler
		centrifuge.ServerJoinHandler
		centrifuge.ServerLeaveHandler
		centrifuge.ErrorHandler
	}
	UnimplementedClientHandlers struct{}

	ISubscriptionHandlers interface {
		centrifuge.PublicationHandler
		centrifuge.JoinHandler
		centrifuge.LeaveHandler
		centrifuge.UnsubscribedHandler
		centrifuge.SubscribingHandler
		centrifuge.SubscribedHandler
		centrifuge.SubscriptionErrorHandler
	}
	UnimplementedSubscriptionHandlers struct{}
)

const (
	BatchState_AtLeastOneAvailable BatchState = iota
	BatchState_PartiallyAvailable
	BatchState_Empty
)

func (s UnimplementedClientHandlers) OnConnectedHandler(centrifuge.ConnectedEvent)                 {}
func (s UnimplementedClientHandlers) OnConnectingHandler(centrifuge.ConnectingEvent)               {}
func (s UnimplementedClientHandlers) OnDisconnectHandler(centrifuge.DisconnectedEvent)             {}
func (s UnimplementedClientHandlers) OnMessageHandler(centrifuge.MessageEvent)                     {}
func (s UnimplementedClientHandlers) OnServerPublicationHandler(centrifuge.ServerPublicationEvent) {}
func (s UnimplementedClientHandlers) OnServerSubscribedHandler(centrifuge.ServerSubscribedEvent)   {}
func (s UnimplementedClientHandlers) OnServerSubscribingHandler(centrifuge.ServerSubscribingEvent) {}
func (s UnimplementedClientHandlers) OnServerUnsubscribedHandler(centrifuge.ServerUnsubscribedEvent) {
}
func (s UnimplementedClientHandlers) OnServerJoinHandler(centrifuge.ServerJoinEvent)   {}
func (s UnimplementedClientHandlers) OnServerLeaveHandler(centrifuge.ServerLeaveEvent) {}
func (s UnimplementedClientHandlers) OnErrorHandler(centrifuge.ErrorEvent)             {}

func (s UnimplementedSubscriptionHandlers) OnPublicationHandler(centrifuge.PublicationEvent) {}
func (s UnimplementedSubscriptionHandlers) OnJoinHandler(centrifuge.JoinEvent)               {}
func (s UnimplementedSubscriptionHandlers) OnLeaveHandler(centrifuge.LeaveEvent)             {}
func (s UnimplementedSubscriptionHandlers) OnUnsubscribedHandler(centrifuge.UnsubscribedEvent) {
}
func (s UnimplementedSubscriptionHandlers) OnSubscribedHandler(centrifuge.SubscribedEvent)   {}
func (s UnimplementedSubscriptionHandlers) OnSubscribingHandler(centrifuge.SubscribingEvent) {}
func (s UnimplementedSubscriptionHandlers) OnSubscriptionErrorHandler(centrifuge.SubscriptionErrorEvent) {
}
