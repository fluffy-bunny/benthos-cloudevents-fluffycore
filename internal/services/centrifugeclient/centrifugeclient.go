package centrifugeclient

import (
	"context"
	"fmt"
	"sync"

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
		mutex  sync.Mutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_centrifuge.ICentrifugeClient = (*service)(nil)

}

func (s *service) Ctor(config *contracts_config.Config) contracts_centrifuge.ICentrifugeClient {
	svc := &service{
		config: config,
	}
	return svc
}

func AddSingletonCentrifugeClient(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_centrifuge.ICentrifugeClient](cb, stemService.Ctor)
}

func (s *service) GetClient() (*centrifuge.Client, error) {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	if s.client == nil {

		client := centrifuge.NewJsonClient(
			s.config.CentrifugeConfig.Endpoint,
			centrifuge.Config{
				// Sending token makes it work with Centrifugo JWT auth (with `secret` HMAC key).
				Token: connToken("49", 0),
			},
		)
		// register for all the events
		client.OnConnected(s.OnConnectedHandler)
		client.OnConnecting(s.OnConnectingHandler)
		client.OnDisconnected(s.OnDisconnectHandler)
		client.OnError(s.OnErrorHandler)
		client.OnJoin(s.OnServerJoinHandler)
		client.OnLeave(s.OnServerLeaveHandler)
		client.OnMessage(s.OnMessageHandler)
		client.OnPublication(s.OnServerPublicationHandler)
		client.OnSubscribed(s.OnServerSubscribedHandler)
		client.OnSubscribing(s.OnServerSubscribingHandler)
		client.OnUnsubscribed(s.OnServerUnsubscribedHandler)
		err := client.Connect()
		if err != nil {
			return nil, err
		}
		s.client = client
	}
	return s.client, nil
}
func (s *service) Dispose(ctx context.Context) error {
	return s.Close(ctx)
}
func (s *service) Close(ctx context.Context) error {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--

	if s.client != nil {
		s.client.Close()
		s.client = nil
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