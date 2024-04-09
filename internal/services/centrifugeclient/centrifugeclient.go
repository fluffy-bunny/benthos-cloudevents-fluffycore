package centrifugeclient

import (
	"context"
	"fmt"
	"sync"

	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/centrifuge"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	jwt "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

type (
	service struct {
		contracts_centrifuge.UnimplementedClientHandlers

		config *contracts_centrifuge.CentrifugeClientConfig
		client *centrifuge.Client
		mutex  sync.Mutex
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_centrifuge.ICentrifugeClient = (*service)(nil)

}

func (s *service) Ctor() contracts_centrifuge.ICentrifugeClient {
	svc := &service{}
	return svc
}

func AddTransientCentrifugeClient(cb di.ContainerBuilder) {
	di.AddTransient[contracts_centrifuge.ICentrifugeClient](cb, stemService.Ctor)
}
func (s *service) Configure(ctx context.Context, config *contracts_centrifuge.CentrifugeClientConfig) {
	log := zerolog.Ctx(ctx).With().Logger()
	if s.config != nil {
		log.Debug().Msg("CentrifugeClientConfig already configured")
		return
	}
	s.config = config
}

func (s *service) GetClient() (*centrifuge.Client, error) {
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//--~--~--~--~--~--~--~--~--~--~-BARBED WIRE-~--~--~--~--~--~--~--~--~--
	if s.client == nil {

		client := centrifuge.NewJsonClient(
			s.config.Endpoint,
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
