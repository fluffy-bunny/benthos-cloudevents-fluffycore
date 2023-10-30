package cloudeventoutput

import (
	benthos_service "github.com/benthosdev/benthos/v4/public/service"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an output to a grpc service.").
	Field(benthos_service.NewStringField("grpc_url")).
	Field(benthos_service.NewStringField("auth_type").Default("none")).
	Field(benthos_service.NewStringField("user_name").Default("none")).
	Field(benthos_service.NewStringField("password").Default("none")).
	Field(benthos_service.NewStringField("api_key").Default("none")).
	Field(benthos_service.NewStringField("api_key_name").Default("x-api-key")).
	Field(benthos_service.NewStringField("client_id").Default("none")).
	Field(benthos_service.NewStringField("client_secret").Default("none")).
	Field(benthos_service.NewStringField("token_endpoint").Default("none")).
	Field(benthos_service.NewStringListField("scopes").Default([]string{""}))

func (s *service) Register() error {
	benthos_service.RegisterOutput("cloudevent_output",
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Output, maxInFlight int, err error) {
			grpcUrl, err := conf.FieldString("grpc_url")
			if err != nil {
				return nil, 0, err
			}
			s.grpcUrl = grpcUrl
			authType, err := conf.FieldString("auth_type")
			if err != nil {
				return nil, 0, err
			}
			switch authType {
			case "basic":
				basicAuthConfig := &basicAuthConfig{}
				s.basicAuthConfig = basicAuthConfig
				value, err := conf.FieldString("user_name")
				if err != nil {
					return nil, 0, err
				}
				basicAuthConfig.UserName = value
				value, err = conf.FieldString("password")
				if err != nil {
					return nil, 0, err
				}
				basicAuthConfig.Password = value
			case "api_key":
				apiKeyConfig := &apiKeyConfig{}
				s.apiKeyConfig = apiKeyConfig
				value, err := conf.FieldString("api_key")
				if err != nil {
					return nil, 0, err
				}
				apiKeyConfig.ApiKey = value
				value, err = conf.FieldString("api_key_name")
				if err != nil {
					return nil, 0, err
				}
				apiKeyConfig.ApiKeyName = value
			case "oauth2":
				oauth2config := &oauth2config{}
				s.oauth2config = oauth2config
				value, err := conf.FieldString("client_secret")
				if err != nil {
					return nil, 0, err
				}
				oauth2config.ClientSecret = value
				value, err = conf.FieldString("client_id")
				if err != nil {
					return nil, 0, err
				}
				oauth2config.ClientId = value
				value, err = conf.FieldString("token_endpoint")
				if err != nil {
					return nil, 0, err
				}
				oauth2config.TokenEndpoint = value
				values, err := conf.FieldStringList("scopes")
				if err != nil {
					return nil, 0, err
				}
				oauth2config.Scopes = values

			}

			return s, 1, nil
		})
	return nil
}
