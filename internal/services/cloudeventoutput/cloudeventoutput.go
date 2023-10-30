package cloudeventoutput

import (
	"context"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	contracts_benthos "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/benthos"
	contracts_cloudeventoutput "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/cloudeventoutput"
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	log "github.com/rs/zerolog/log"
	cc "golang.org/x/oauth2/clientcredentials"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	metadata "google.golang.org/grpc/metadata"
	protojson "google.golang.org/protobuf/encoding/protojson"
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
		grpcUrl         string
		oauth2config    *oauth2config
		apiKeyConfig    *apiKeyConfig
		basicAuthConfig *basicAuthConfig
		// authType: none, oauth2, api_key
		authType                  string
		cloudEventProcessorClient proto_cloudeventprocessor.CloudEventProcessorClient
	}
)

var stemService = &service{}

func (s *service) Ctor() *service {
	return &service{}
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

func (s *service) Connect(ctx context.Context) error {
	// https://github.com/johanbrandhorst/grpc-auth-example
	/*
	   using middleware to add the authorization header but concept is the same.
	*/
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if s.basicAuthConfig != nil {
		outboundAuthMiddleware := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Basic "+s.basicAuthConfig.UserName+":"+s.basicAuthConfig.Password)
			return invoker(newCtx, method, req, reply, cc, opts...)
		}
		// add the authorization header
		do := grpc.WithUnaryInterceptor(outboundAuthMiddleware)
		dialOptions = append(dialOptions, do)
	}
	if s.apiKeyConfig != nil {
		outboundAuthMiddleware := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			newCtx := metadata.AppendToOutgoingContext(ctx, s.apiKeyConfig.ApiKeyName, s.apiKeyConfig.ApiKey)
			return invoker(newCtx, method, req, reply, cc, opts...)
		}
		// add the authorization header
		do := grpc.WithUnaryInterceptor(outboundAuthMiddleware)
		dialOptions = append(dialOptions, do)
	}
	if s.oauth2config != nil {
		ccConfig := &cc.Config{
			ClientID:     s.oauth2config.ClientId,
			ClientSecret: s.oauth2config.ClientSecret,
			TokenURL:     s.oauth2config.TokenEndpoint,
			Scopes:       s.oauth2config.Scopes,
		}
		tokenSource := ccConfig.TokenSource(context.Background())
		_, err := tokenSource.Token()
		if err != nil {
			log.Error().Err(err).Msg("failed to get token")
			return err
		}
		// grpc clients can have middleware, in this case we make sure the token is added to the request
		outboundAuthMiddleware := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			// add authorization header
			token, err := tokenSource.Token() // get the token from the source, it may have expired so it will get a new one.
			if err != nil {
				log.Error().Err(err).Msg("Failed to get token")
				return err
			}
			newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
			return invoker(newCtx, method, req, reply, cc, opts...)
		}
		// add the authorization header
		do := grpc.WithUnaryInterceptor(outboundAuthMiddleware)
		dialOptions = append(dialOptions, do)

	}
	conn, err := grpc.Dial(s.grpcUrl, dialOptions...)
	if err != nil {
		log.Error().Err(err).Msg("Dial failed")
		return err
	}

	s.cloudEventProcessorClient = proto_cloudeventprocessor.NewCloudEventProcessorClient(conn)

	return nil
}
func (s *service) Write(ctx context.Context, message *benthos_service.Message) error {
	content, err := message.AsBytes()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert message to bytes")
		return err
	}
	ceBatch := &cloudevents.CloudEventBatch{}
	err = protojson.Unmarshal(content, ceBatch)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal CloudEventBatch")
		return err
	}
	log.Info().Interface("ceBatch", ceBatch).Msg("ceBatch")
	return nil
}
func (s *service) Close(ctx context.Context) error {
	return nil
}
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
