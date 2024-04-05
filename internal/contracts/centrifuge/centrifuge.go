package centrifuge

import (
	centrifuge "github.com/centrifugal/centrifuge-go"
	fluffycore_contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
)

type (
	ICentrifugeClient interface {
		fluffycore_contracts_common.ICloser
		fluffycore_contracts_common.IDispose

		GetClient() *centrifuge.Client
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
func (s UnimplementedSubscriptionHandlers) OnJoinHandler(centrifuge.JoinEvent)          {}
func (s UnimplementedSubscriptionHandlers) OnLeaveHandler(centrifuge.LeaveEvent)        {}
func (s UnimplementedSubscriptionHandlers) OnUnsubscribedHandler(centrifuge.UnsubscribedEvent) {
}
func (s UnimplementedSubscriptionHandlers) OnSubscribedHandler(centrifuge.SubscribedEvent)   {}
func (s UnimplementedSubscriptionHandlers) OnSubscribingHandler(centrifuge.SubscribingEvent) {}
func (s UnimplementedSubscriptionHandlers) OnSubscriptionErrorHandler(centrifuge.SubscriptionErrorEvent) {
}
