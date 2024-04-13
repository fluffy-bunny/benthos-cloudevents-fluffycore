package exampleSingletonTokenAccessor

/*
Example of a service that implements the
*/
import (
	"fmt"

	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/contracts/centrifuge"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	jwt "github.com/golang-jwt/jwt"
)

type (
	service struct {
	}
)

var stemService = (*service)(nil)

func init() {
	var _ contracts_centrifuge.ISingletonCentrifugeTokenAccessor = (*service)(nil)
}

func (s *service) Ctor() contracts_centrifuge.ISingletonCentrifugeTokenAccessor {
	return &service{}
}

func AddSingletonCentrifugeTokenAccessor(cb di.ContainerBuilder) {
	di.AddSingleton[contracts_centrifuge.ISingletonCentrifugeTokenAccessor](cb, stemService.Ctor)
}
func (s *service) GetToken(evt centrifuge.ConnectionTokenEvent) (string, error) {
	return ExampleConnToken("49", 0), nil
}

const ExampleTokenHmacSecret = "my_secret"

func ExampleConnToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ExampleTokenHmacSecret))
	if err != nil {
		panic(err)
	}
	fmt.Println("token: ", t)
	return t
}
