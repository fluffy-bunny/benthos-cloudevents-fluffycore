/*
 Copyright 2021 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka_franz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	proto_cloudevents "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudevents"
	"github.com/gogo/status"
	zLog "github.com/rs/zerolog/log"
	kgo "github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/grpc/codes"
)

type withMessageKey struct{}

// WithMessageKey allows to set the key used when sending the producer message
func WithMessageKey(ctx context.Context, key []byte) context.Context {
	return context.WithValue(ctx, withMessageKey{}, key)
}

// CloudEventToKafkaMessage converts a grpc CloudEvent to a kafka.Message
func CloudEventToKafkaMessage(ctx context.Context, ce *proto_cloudevents.CloudEvent) (*kgo.Record, error) {
	e := cloudevents.NewEvent()
	e.SetID(ce.Id)
	e.SetType(ce.Type)
	e.SetSource(ce.Source)
	attrib, ok := ce.Attributes["time"]
	if !ok {
		e.SetTime(time.Now())
	} else {
		timeStamp := attrib.GetCeTimestamp()
		e.SetTime(timeStamp.AsTime())
	}
	// hardcoded to json for now
	var obj interface{}
	textData := ce.GetTextData()
	err := json.Unmarshal([]byte(textData), &obj)
	if err != nil {
		return nil, err
	}
	err = e.SetData(cloudevents.ApplicationJSON, obj)
	if err != nil {
		return nil, err
	}
	return MakeKafkaMessage(ctx, e)
}

// MakeKafkaMessage from cloudEvent sdk-go Event
func MakeKafkaMessage(ctx context.Context, event cloudevents.Event) (*kgo.Record, error) {
	log := zLog.With().Str("source", event.Source()).Logger()
	m := (*binding.EventMessage)(&event)
	if m == nil {
		err := status.Error(codes.InvalidArgument, "event can't be typcast to EventMessage")
		log.Error().Err(err).Send()
		return nil, err
	}
	return makeKafkaMessage2(ctx, (*binding.EventMessage)(&event))
}
func makeKafkaMessage2(ctx context.Context, m binding.Message, transformers ...binding.Transformer) (*kgo.Record, error) {
	var err error
	defer m.Finish(err)

	kafkaMessage := kgo.Record{}

	if k := ctx.Value(withMessageKey{}); k != nil {
		kafkaMessage.Key = []byte(fmt.Sprintf("%v", k.(interface{})))
	}

	if err = WriteProducerMessage(ctx, m, &kafkaMessage, transformers...); err != nil {
		return nil, err
	}
	return &kafkaMessage, nil
}
