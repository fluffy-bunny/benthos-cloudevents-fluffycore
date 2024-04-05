package config

import (
	fluffycore_contracts_ddprofiler "github.com/fluffy-bunny/fluffycore/contracts/ddprofiler"
)

type (
	KafkaConfig struct {
		Seeds []string `json:"seeds"`
		Group string   `json:"group"`
		Topic string   `json:"topic"`
	}
	InitialConfig struct {
		ConfigFiles ConfigFiles `json:"configFiles"`
	}
	ConfigFiles struct {
		RootPath string `json:"rootPath"`
	}

	CentrifugeConfig struct {
		Endpoint string `json:"endpoint"`
	}
	Config struct {
		DDProfilerConfig      *fluffycore_contracts_ddprofiler.Config `json:"ddProfilerConfig"`
		KafkaDeadLetterConfig *KafkaConfig                            `json:"kafkaDeadLetterConfig"`
		EnableKafkaSASL       bool                                    `json:"enableKafkaSASL"`
		ConfigFiles           *ConfigFiles                            `json:"configFiles"`
		CentrifugeConfig      *CentrifugeConfig                       `json:"centrifugeConfig"`
	}
)

// ConfigDefaultJSON default json
var ConfigDefaultJSON = []byte(`
{
	"enableKafkaSASL": false,
	"ddProfilerConfig": {
		"enabled": false,
		"serviceName": "in-environment",
		"applicationEnvironment": "in-environment",
		"version": "1.0.0"
	},
	"kafkaDeadLetterConfig": {
		"seeds": ["localhost:9093"],
		"group": "$Default",
		"topic": "cloudevents-core-deadletter"
	},
	"configFiles": {
		"rootPath": ""
	},
	"centrifugeConfig": {
		"endpoint": "ws://localhost:8000/connection/websocket"
	}

  }
`)
