package kafkaclient

import (
	contracts_config "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/cmd/processor/internal/contracts/config"
	di "github.com/fluffy-bunny/fluffy-dozm-di"
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

func AddSingletonKafkaClient(builder di.ContainerBuilder) {
	di.AddSingleton[*kgo.Client](builder, func(config *contracts_config.Config) (*kgo.Client, error) {
		cl, err := kgo.NewClient(
			kgo.SeedBrokers(config.KafkaConfig.Seeds...),
			kgo.ConsumerGroup(config.KafkaConfig.Group),
			kgo.ConsumeTopics(config.KafkaConfig.Topic),
		)
		if err != nil {
			return nil, err
		}
		return cl, nil
	})

}
