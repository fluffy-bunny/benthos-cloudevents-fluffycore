package kafkaclient

import (
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/config"
	contracts_kafkaclient "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/kafkaclient"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

type (
	service struct {
		config *contracts_config.KafkaConfig
		client *kgo.Client
	}
)

func init() {
	var _ contracts_kafkaclient.IKafkaClient = (*service)(nil)
}

func AddSingletonKafkaDeadLetterClient(builder di.ContainerBuilder) {
	di.AddSingleton[contracts_kafkaclient.IDeadLetterClient](builder, func(config *contracts_config.Config) (contracts_kafkaclient.IDeadLetterClient, error) {
		cl, err := kgo.NewClient(
			kgo.SeedBrokers(config.KafkaDeadLetterConfig.Seeds...),
			kgo.DefaultProduceTopic(config.KafkaDeadLetterConfig.Topic),
		)
		if err != nil {
			return nil, err
		}
		return &service{
			config: config.KafkaDeadLetterConfig,
			client: cl,
		}, nil

	})
}

func (s *service) GetClient() *kgo.Client {
	return s.client
}
