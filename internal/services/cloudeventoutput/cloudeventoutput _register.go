package cloudeventoutput

import (
	"strings"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
)

const (
	// Output
	OutputName            = "cloudevents_grpc"
	FieldName_gRPCUrl     = "grpc_url"
	FieldName_channel     = "channel"
	FieldName_maxInFlight = "max_in_flight"
	// Auth
	FieldName_auth = "auth"
	// Auth.apikey
	FieldName_apikey = "apikey"
	FieldName_name   = "name"
	FieldName_value  = "value"
	// Auth.basic
	FieldName_basic    = "basic"
	FieldName_userName = "user_name"
	FieldName_password = "password"
	// Auth.oauth2
	FieldName_oauth2        = "oauth2"
	FieldName_clientId      = "client_id"
	FieldName_clientSecret  = "client_secret"
	FieldName_tokenEndpoint = "token_endpoint"
	FieldName_scopes        = "scopes"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an output to a grpc service.").
	Field(benthos_service.NewStringField(FieldName_gRPCUrl)).
	Field(benthos_service.NewStringField(FieldName_channel)).
	Field(benthos_service.NewInputMaxInFlightField()).
	Field(NewAuthConfig().Optional())

func NewBasicAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField(FieldName_basic,
		benthos_service.NewStringField(FieldName_userName),
		benthos_service.NewStringField(FieldName_password),
	)
}
func NewOAuth2AuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField(FieldName_oauth2,
		benthos_service.NewStringField(FieldName_clientId),
		benthos_service.NewStringField(FieldName_clientSecret),
		benthos_service.NewStringField(FieldName_tokenEndpoint),
		benthos_service.NewStringListField(FieldName_scopes).Default([]string{""}),
	)
}
func NewAPIKeyAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField(FieldName_apikey,
		benthos_service.NewStringField(FieldName_name),
		benthos_service.NewStringField(FieldName_value),
	)
}
func NewAuthConfig() *benthos_service.ConfigField {
	return benthos_service.NewObjectField(FieldName_auth,
		NewBasicAuthConfig().Optional(),
		NewOAuth2AuthConfig().Optional(),
		NewAPIKeyAuthConfig().Optional(),
	)
}

func (s *service) Register() error {
	benthos_service.RegisterOutput(OutputName,
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Output, maxInFlight int, err error) {
			grpcUrl, err := conf.FieldString(FieldName_gRPCUrl)
			if err != nil {
				return nil, 0, err
			}
			s.grpcUrl = grpcUrl

			channel, err := conf.FieldString(FieldName_channel)
			if err != nil {
				return nil, 0, err
			}
			s.channel = channel

			maxInFlight, err = conf.FieldInt(FieldName_maxInFlight)
			if err != nil {
				return nil, 0, err
			}

			auth, err := conf.FieldAnyMap(FieldName_auth)
			if err != nil {
				// optional
				if !strings.Contains(err.Error(), "not found") {
					return nil, 0, err
				}
			}
			if len(auth) > 0 {
				ok := false
				authTypes := []string{
					FieldName_oauth2,
					FieldName_basic,
					FieldName_apikey,
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
					case FieldName_apikey:
						apiKeyConfig := &apiKeyConfig{}
						s.apiKeyConfig = apiKeyConfig
						value, err := parsedConfig.FieldString(FieldName_value)
						if err != nil {
							return nil, 0, err
						}
						apiKeyConfig.ApiKey = value
						name, err := parsedConfig.FieldString(FieldName_name)
						if err != nil {
							return nil, 0, err
						}
						apiKeyConfig.ApiKeyName = name
						ok = true
					case FieldName_basic:
						basicAuthConfig := &basicAuthConfig{}
						s.basicAuthConfig = basicAuthConfig
						value, err := parsedConfig.FieldString(FieldName_userName)
						if err != nil {
							return nil, 0, err
						}
						basicAuthConfig.UserName = value
						value, err = parsedConfig.FieldString(FieldName_password)
						if err != nil {
							return nil, 0, err
						}
						basicAuthConfig.Password = value
						ok = true
					case FieldName_oauth2:
						oauth2config := &oauth2config{}
						s.oauth2config = oauth2config
						value, err := parsedConfig.FieldString(FieldName_clientSecret)
						if err != nil {
							return nil, 0, err
						}
						oauth2config.ClientSecret = value
						value, err = parsedConfig.FieldString(FieldName_clientId)
						if err != nil {
							return nil, 0, err
						}
						oauth2config.ClientId = value
						value, err = parsedConfig.FieldString(FieldName_tokenEndpoint)
						if err != nil {
							return nil, 0, err
						}
						oauth2config.TokenEndpoint = value
						values, err := parsedConfig.FieldStringList(FieldName_scopes)
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
