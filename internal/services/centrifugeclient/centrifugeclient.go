package centrifugeclient

import (
	"context"
	"fmt"

	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/centrifuge"
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	jwt "github.com/golang-jwt/jwt"
)

type (
	service struct {
		contracts_centrifuge.UnimplementedClientHandlers

		config *contracts_config.Config
		client *centrifuge.Client
	}
)

var stemService = (*service)(nil)

func (s *service) Ctor(config *contracts_config.Config) contracts_centrifuge.ICentrifugeClient {

	client := centrifuge.NewJsonClient(
		config.CentrifugeConfig.Endpoint,
		centrifuge.Config{
			// Sending token makes it work with Centrifugo JWT auth (with `secret` HMAC key).
			Token: connToken("49", 0),
		},
	)

	svc := &service{
		config: config,
		client: client,
	}

	// register for all the events
	client.OnConnected(svc.OnConnectedHandler)
	client.OnConnecting(svc.OnConnectingHandler)
	client.OnDisconnected(svc.OnDisconnectHandler)
	client.OnError(svc.OnErrorHandler)
	client.OnJoin(svc.OnServerJoinHandler)
	client.OnLeave(svc.OnServerLeaveHandler)
	client.OnMessage(svc.OnMessageHandler)
	client.OnPublication(svc.OnServerPublicationHandler)
	client.OnSubscribed(svc.OnServerSubscribedHandler)
	client.OnSubscribing(svc.OnServerSubscribingHandler)
	client.OnUnsubscribed(svc.OnServerUnsubscribedHandler)

	return svc

}
func init() {
	var _ contracts_centrifuge.ICentrifugeClient = (*service)(nil)

}
func AddSingletonCentrifugeClient(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_centrifuge.ICentrifugeClient](cb, stemService.Ctor)
}

func (s *service) GetClient() *centrifuge.Client {
	return s.client
}
func (s *service) Dispose(ctx context.Context) error {
	return s.Close(ctx)
}
func (s *service) Close(ctx context.Context) error {
	if s.client != nil {
		s.client.Close()
	}
	return nil
}

const exampleTokenHmacSecret = "my_secret"

func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(exampleTokenHmacSecret))
	if err != nil {
		panic(err)
	}
	fmt.Println("token: ", t)
	return t
}
