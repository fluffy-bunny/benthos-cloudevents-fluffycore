package config

import (
	fluffycore_contracts_config "github.com/fluffy-bunny/fluffycore/contracts/config"
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
)

type (
	JWTValidators struct {
		Issuers  []string `json:"issuers" mapstructure:"ISSUERS"`
		JWKSURLS []string `json:"jwksUrls" mapstructure:"JWKS_URLS"`
	}
	KafkaConfig struct {
		Seeds []string `json:"seeds" mapstructure:"SEEDS"`
		Group string   `json:"group" mapstructure:"GROUP"`
		Topic string   `json:"topic" mapstructure:"TOPIC"`
	}
	Config struct {
		fluffycore_contracts_config.CoreConfig `mapstructure:",squash"`
		CustomString                           string                                  `json:"CUSTOM_STRING" mapstructure:"CUSTOM_STRING"`
		SomeSecret                             string                                  `json:"SOME_SECRET" mapstructure:"SOME_SECRET" redact:"true"`
		OAuth2Port                             int                                     `json:"oauth2Port"  mapstructure:"OAUTH2_PORT"`
		JWTValidators                          JWTValidators                           `json:"jwtValidators" mapstructure:"JWT_VALIDATORS"`
		DDProfilerConfig                       *fluffycore_contracts_ddprofiler.Config `json:"ddProfilerConfig" mapstructure:"DD_PROFILER_CONFIG"`
		KafkaConfig                            KafkaConfig                             `json:"kafkaConfig" mapstructure:"KAFKA_CONFIG"`
	}
)

// ConfigDefaultJSON default json
var ConfigDefaultJSON = []byte(`
{
	"APPLICATION_NAME": "in-environment",
	"APPLICATION_ENVIRONMENT": "in-environment",
	"PRETTY_LOG": false,
	"LOG_LEVEL": "info",
	"PORT": 50051,
	"REST_PORT": 50052,
	"OAUTH2_PORT": 50053,
	"CUSTOM_STRING": "some default value",
	"SOME_SECRET": "password",
	"GRPC_GATEWAY_ENABLED": true,
	"JWT_VALIDATORS": {},
	"DD_PROFILER_CONFIG": {
		"ENABLED": false,
		"SERVICE_NAME": "in-environment",
		"APPLICATION_ENVIRONMENT": "in-environment",
		"VERSION": "1.0.0"
	},
	"KAFKA_CONFIG": {
		"SEEDS": ["host.docker.internal:9093"],
		"GROUP": "$Default",
		"TOPIC": "cloudevents-core"
	}

  }
`)
