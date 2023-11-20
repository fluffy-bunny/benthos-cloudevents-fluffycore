package config

import (
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
)

type (
	KafkaConfig struct {
		Seeds []string `json:"seeds" mapstructure:"SEEDS"`
		Group string   `json:"group" mapstructure:"GROUP"`
		Topic string   `json:"topic" mapstructure:"TOPIC"`
	}
	Config struct {
		DDProfilerConfig      *fluffycore_contracts_ddprofiler.Config `json:"ddProfilerConfig" mapstructure:"DD_PROFILER_CONFIG"`
		KafkaDeadLetterConfig *KafkaConfig                            `json:"kafkaDeadLetterConfig" mapstructure:"KAFKA_DEAD_LETTER_CONFIG"`
	}
)

// ConfigDefaultJSON default json
var ConfigDefaultJSON = []byte(`
{
	"DD_PROFILER_CONFIG": {
		"ENABLED": false,
		"SERVICE_NAME": "in-environment",
		"APPLICATION_ENVIRONMENT": "in-environment",
		"VERSION": "1.0.0"
	},
	"KAFKA_DEAD_LETTER_CONFIG": {
		"SEEDS": ["localhost:9093"],
		"GROUP": "$Default",
		"TOPIC": "cloudevents-core-deadletter"
	}

  }
`)
