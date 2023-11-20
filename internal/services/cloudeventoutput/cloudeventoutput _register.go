package cloudeventoutput

import (
	"strings"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an output to a grpc service.").
	Field(benthos_service.NewStringField("grpc_url")).
	Field(benthos_service.NewInputMaxInFlightField()).
	Field(NewAuthConfig().Optional()).
	Field(NewKafkaDeadLetterConfig().Optional())

func NewKafkaDeadLetterConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("kafka_franz_dead_letter",
		benthos_service.NewStringListField("seed_brokers"),
		benthos_service.NewStringField("topic"),
		benthos_service.NewStringField("key").Optional(),
		NewSASLConfig().Optional(),
	)
}
func NewSASLConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("sasl",
		benthos_service.NewStringField("mechanism"),
		benthos_service.NewStringField("username"),
		benthos_service.NewStringField("password"),
	)
}
func NewBasicAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("basic",
		benthos_service.NewStringField("user_name"),
		benthos_service.NewStringField("password"),
	)
}
func NewOAuth2AuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("oauth2",
		benthos_service.NewStringField("client_id"),
		benthos_service.NewStringField("client_secret"),
		benthos_service.NewStringField("token_endpoint"),
		benthos_service.NewStringListField("scopes").Default([]string{""}),
	)
}
func NewAPIKeyAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("apikey",
		benthos_service.NewStringField("name"),
		benthos_service.NewStringField("value"),
	)
}
func NewAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField("auth",
		NewBasicAuthConfig().Optional(),
		NewOAuth2AuthConfig().Optional(),
		NewAPIKeyAuthConfig().Optional(),
	)
}

type (
	kafkaFranzDeadLetterConfig struct {
		SeedBrokers []string
		Topic       string
		Key         string
		SASL        *SASLConfig
	}
	SASLConfig struct {
		Mechanism string
		Username  string
		Password  string
	}
)

func (s *service) Register() error {
	benthos_service.RegisterOutput("cloudevents_grpc",
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Output, maxInFlight int, err error) {
			grpcUrl, err := conf.FieldString("grpc_url")
			if err != nil {
				return nil, 0, err
			}
			s.grpcUrl = grpcUrl
			maxInFlight, err = conf.FieldInt("max_in_flight")
			if err != nil {
				return nil, 0, err
			}
			loadBasicDeadLetterConfig := func() (*kafkaFranzDeadLetterConfig, error) {
				seedBrokers, err := conf.FieldStringList("kafka_franz_dead_letter", "seed_brokers")
				if err != nil {
					return nil, err
				}
				topic, err := conf.FieldString("kafka_franz_dead_letter", "topic")
				if err != nil {
					return nil, err
				}
				key, err := conf.FieldString("kafka_franz_dead_letter", "key")
				if err != nil {
					key = ""
				}
				return &kafkaFranzDeadLetterConfig{
					SeedBrokers: seedBrokers,
					Topic:       topic,
					Key:         key,
				}, nil
			}
			if conf.Contains("kafka_franz_dead_letter") {
				kfdlConfig, err := loadBasicDeadLetterConfig()
				if err != nil {
					return nil, 0, err
				}
				if kfdlConfig != nil {
					s.kafkaFranzDeadLetter = kfdlConfig
					// next get the SASL config
					kafkaFranzDeadLetterConfig, err := conf.FieldAnyMap("kafka_franz_dead_letter")
					if err == nil {
						pc, ok := kafkaFranzDeadLetterConfig["sasl"]
						if ok {
							mechanism, err := pc.FieldString("mechanism")
							if err != nil {
								return nil, 0, err
							}
							username, err := pc.FieldString("username")
							if err != nil {
								return nil, 0, err
							}
							password, err := pc.FieldString("password")
							if err != nil {
								return nil, 0, err
							}
							saslConfig := &SASLConfig{
								Mechanism: mechanism,
								Username:  username,
								Password:  password,
							}
							kfdlConfig.SASL = saslConfig
						}
					}
				}
			}
			auth, err := conf.FieldAnyMap("auth")
			if err != nil {
				// optional
				if !strings.Contains(err.Error(), "not found") {
					return nil, 0, err
				}
			}
			if len(auth) > 0 {
				ok := false
				authTypes := []string{
					"oauth2",
					"basic",
					"apikey",
				}
				for _, v := range authTypes {
					if ok {
						break
					}
					parsedConfig, ok2 := auth[v]
					if !ok2 {
						continue
					}

					switch v {
					case "apikey":
						apiKeyConfig := &apiKeyConfig{}
						s.apiKeyConfig = apiKeyConfig
						value, err := parsedConfig.FieldString("value")
						if err != nil {
							return nil, 0, err
						}
						apiKeyConfig.ApiKey = value
						name, err := parsedConfig.FieldString("name")
						if err != nil {
							return nil, 0, err
						}
						apiKeyConfig.ApiKeyName = name
						ok = true
					case "basic":
						basicAuthConfig := &basicAuthConfig{}
						s.basicAuthConfig = basicAuthConfig
						value, err := parsedConfig.FieldString("user_name")
						if err != nil {
							return nil, 0, err
						}
						basicAuthConfig.UserName = value
						value, err = parsedConfig.FieldString("password")
						if err != nil {
							return nil, 0, err
						}
						basicAuthConfig.Password = value
						ok = true
					case "oauth2":
						oauth2config := &oauth2config{}
						s.oauth2config = oauth2config
						value, err := parsedConfig.FieldString("client_secret")
						if err != nil {
							return nil, 0, err
						}
						oauth2config.ClientSecret = value
						value, err = parsedConfig.FieldString("client_id")
						if err != nil {
							return nil, 0, err
						}
						oauth2config.ClientId = value
						value, err = parsedConfig.FieldString("token_endpoint")
						if err != nil {
							return nil, 0, err
						}
						oauth2config.TokenEndpoint = value
						values, err := parsedConfig.FieldStringList("scopes")
						if err != nil {
							return nil, 0, err
						}
						oauth2config.Scopes = values
						ok = true
					}
				}
				if !ok {
					return nil, maxInFlight, err
				}
			}

			return s, maxInFlight, nil
		})
	return nil
}
