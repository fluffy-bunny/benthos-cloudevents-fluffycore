package runtime

import (
	"context"
	"encoding/json"

	internal_auth "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/auth"
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	services_health "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/services/health"
	services_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/services/kafkaclient"
	services_kafkacloudeventservice "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/services/kafkacloudeventservice"
	services_processor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/services/processor"
	internal_version "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/version"
	pkg_services_centrifuge_batcher "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/services/centrifuge/batcher"
	pkg_services_centrifuge_client "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/services/centrifuge/client"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
	fluffycore_contracts_middleware "github.com/fluffy-bunny/fluffycore/contracts/middleware"
	fluffycore_contracts_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/contracts/middleware/auth/jwt"
	fluffycore_contracts_runtime "github.com/fluffy-bunny/fluffycore/contracts/runtime"
	fluffycore_middleware_auth_jwt "github.com/fluffy-bunny/fluffycore/middleware/auth/jwt"
	fluffycore_middleware_claimsprincipal "github.com/fluffy-bunny/fluffycore/middleware/claimsprincipal"
	fluffycore_middleware_correlation "github.com/fluffy-bunny/fluffycore/middleware/correlation"
	fluffycore_middleware_dicontext "github.com/fluffy-bunny/fluffycore/middleware/dicontext"
	fluffycore_middleware_logging "github.com/fluffy-bunny/fluffycore/middleware/logging"
	core_runtime "github.com/fluffy-bunny/fluffycore/runtime"
	fluffycore_services_ddprofiler "github.com/fluffy-bunny/fluffycore/services/ddprofiler"
	fluffycore_utils_redact "github.com/fluffy-bunny/fluffycore/utils/redact"
	status "github.com/gogo/status"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	codes "google.golang.org/grpc/codes"
)

type (
	startup struct {
		fluffycore_contracts_runtime.UnimplementedStartup
		RootContainer di.Container
		config        *contracts_config.Config
		configOptions *fluffycore_contracts_runtime.ConfigOptions
		ddProfiler    fluffycore_contracts_ddprofiler.IDataDogProfiler
	}
)

func NewStartup() fluffycore_contracts_runtime.IStartup {
	return &startup{}
}
func (s *startup) ConfigureServices(ctx context.Context, builder di.ContainerBuilder) {
	log := zerolog.Ctx(ctx).With().Str("startup:method", "ConfigureServices").Logger()
	dst, err := fluffycore_utils_redact.CloneAndRedact(s.configOptions.Destination)
	if err != nil {
		panic(err)
	}
	log.Info().Interface("config", dst).Msg("config")
	config := s.configOptions.Destination.(*contracts_config.Config)
	config.DDProfilerConfig.ApplicationEnvironment = config.ApplicationEnvironment
	config.DDProfilerConfig.ServiceName = config.ApplicationName
	config.DDProfilerConfig.Version = internal_version.Version()

	// Add our configs to the di as simple objects
	//--------------------------------------
	di.AddInstance[*fluffycore_contracts_ddprofiler.Config](builder, config.DDProfilerConfig)
	di.AddInstance[*contracts_config.Config](builder, config)
	fluffycore_services_ddprofiler.AddSingletonIProfiler(builder)

	// Add Health Service
	//--------------------------------------
	services_health.AddHealthService(builder)

	// Add Authorization
	//--------------------------------------
	issuerConfigs := &fluffycore_contracts_middleware_auth_jwt.IssuerConfigs{}
	for idx := range s.config.JWTValidators.Issuers {
		issuerConfigs.IssuerConfigs = append(issuerConfigs.IssuerConfigs,
			&fluffycore_contracts_middleware_auth_jwt.IssuerConfig{
				OAuth2Config: &fluffycore_contracts_middleware_auth_jwt.OAuth2Config{
					Issuer:  s.config.JWTValidators.Issuers[idx],
					JWKSUrl: s.config.JWTValidators.JWKSURLS[idx],
				},
			})
	}
	fluffycore_middleware_auth_jwt.AddValidators(builder, issuerConfigs)

	// Add CloudEventProcessor Service
	//--------------------------------------
	services_processor.AddCloudEventProcessorServer(builder)
	services_kafkacloudeventservice.AddKafkaCloudEventServiceServer(builder)
	services_kafkaclient.AddSingletonKafkaDeadLetterClient(builder)
	services_kafkaclient.AddSingletonKafkaPublishingClient(builder)
	pkg_services_centrifuge_batcher.AddTransientCentrifugeStreamBatcher(builder)
	pkg_services_centrifuge_client.AddTransientCentrifugeClient(builder)

}
func (s *startup) Configure(ctx context.Context, rootContainer di.Container, unaryServerInterceptorBuilder fluffycore_contracts_middleware.IUnaryServerInterceptorBuilder, streamServerInterceptorBuilder fluffycore_contracts_middleware.IStreamServerInterceptorBuilder) {
	log := zerolog.Ctx(ctx).With().Str("method", "Configure").Logger()

	// puts a zerlog logger into the request context
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_logging.EnsureContextLoggingUnaryServerInterceptor())

	// log correlation and spans
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_correlation.EnsureCorrelationIDUnaryServerInterceptor())
	// dicontext is responsible of create a scoped context for each request.
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_dicontext.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_dicontext.UnaryServerInterceptor(rootContainer))
	// auth
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_auth_jwt.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_auth_jwt.UnaryServerInterceptor(rootContainer))

	// Here the gating happens
	grpcEntrypointClaimsMap := internal_auth.BuildGrpcEntrypointPermissionsClaimsMap()
	// claims principal
	log.Info().Msg("adding unaryServerInterceptorBuilder: fluffycore_middleware_claimsprincipal.UnaryServerInterceptor")
	unaryServerInterceptorBuilder.Use(fluffycore_middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMapWithZeroTrustV2(grpcEntrypointClaimsMap))

	// last is the recovery middleware
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	unaryServerInterceptorBuilder.Use(grpc_recovery.UnaryServerInterceptor(opts...))
}
func (s *startup) GetConfigOptions() *fluffycore_contracts_runtime.ConfigOptions {
	log := log.With().Caller().Str("method", "GetConfigOptions").Logger()
	// here we load a config file and merge it over the default.
	initialConfigOptions := &fluffycore_contracts_runtime.ConfigOptions{
		Destination: &contracts_config.InitialConfig{},
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	err := core_runtime.LoadConfig(initialConfigOptions)
	if err != nil {
		panic(err)
	}

	defaultConfig := &contracts_config.Config{}
	err = json.Unmarshal([]byte(contracts_config.ConfigDefaultJSON), defaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal ConfigDefaultJSON")
	}
	log.Info().Interface("defaultConfig", defaultConfig).Msg("config after merge")
	s.config = &contracts_config.Config{}
	s.configOptions = &fluffycore_contracts_runtime.ConfigOptions{
		Destination: s.config,
		RootConfig:  contracts_config.ConfigDefaultJSON,
	}
	return s.configOptions
}

func (s *startup) SetRootContainer(container di.Container) {
	s.RootContainer = container
}
