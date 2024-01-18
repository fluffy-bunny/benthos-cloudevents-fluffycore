package kafkaclient

import (
	kgo "github.com/twmb/franz-go/pkg/kgo"
)

type (
	IKafkaClient interface {
		GetClient() *kgo.Client
	}
	IPublishingClient interface {
		IKafkaClient
	}
	IDeadLetterClient interface {
		IKafkaClient
	}
)
