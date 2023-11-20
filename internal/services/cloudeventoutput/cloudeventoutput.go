package cloudeventoutput

import (
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_cloudeventoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/cloudeventoutput"
	contracts_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/kafkaclient"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
)

type (
	oauth2config struct {
		ClientId      string
		ClientSecret  string
		TokenEndpoint string
		Scopes        []string
	}
	apiKeyConfig struct {
		ApiKey     string
		ApiKeyName string
	}
	basicAuthConfig struct {
		UserName string
		Password string
	}
	service struct {
		// grpcUrl: i.e. grpc://localhost:5001
		grpcUrl                   string
		oauth2config              *oauth2config
		apiKeyConfig              *apiKeyConfig
		basicAuthConfig           *basicAuthConfig
		cloudEventProcessorClient proto_cloudeventprocessor.CloudEventProcessorClient
		deadLetterClient          contracts_kafkaclient.IDeadLetterClient
	}
)

var stemService = &service{}

func (s *service) Ctor(deadLetterClient contracts_kafkaclient.IDeadLetterClient) *service {
	return &service{
		deadLetterClient: deadLetterClient,
	}
}
func init() {
	var _ contracts_cloudeventoutput.ICloudEventOutput = (*service)(nil)
	var _ contracts_benthos.IBenthosRegistration = (*service)(nil)
}
func AddSingletonCloudEventOutput(cb di.ContainerBuilder) {
	di.AddSingleton[*service](cb, stemService.Ctor,
		contracts_benthos.TypeIBenthosRegistration,
		contracts_cloudeventoutput.TypeICloudEventOutput)
}
