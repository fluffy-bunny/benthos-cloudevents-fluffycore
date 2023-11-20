package cloudeventoutput

import (
	"context"

	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	log "github.com/rs/zerolog/log"
	cc "golang.org/x/oauth2/clientcredentials"
	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	metadata "google.golang.org/grpc/metadata"
)

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
