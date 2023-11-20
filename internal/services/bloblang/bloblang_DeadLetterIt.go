package bloblang

import (
	"context"
	"sync"
	"time"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

/*
	{
	    "content": "ZA==",
	    "headers": {
	        "kafka_key": "",
	        "kafka_lag": 0,
	        "kafka_offset": 64,
	        "kafka_partition": 0,
	        "kafka_timestamp_unix": 1700508875,
	        "kafka_tombstone_message": false,
	        "kafka_topic": "cloudevents-core"
	    },
	    "type": "error"
	}
*/
func (s *service) DeadLetterIt(args *bloblang.ParsedParams) (bloblang.Function, error) {
	headers, err := args.Get("headers")
	if err != nil {
		return nil, err
	}
	// convert headers into a map
	headersMap, ok := headers.(map[string]interface{})
	if !ok {
		return nil, err
	}
	// trim the kafka_ prefix
	headersList := make([]string, 0, len(headersMap))
	for k := range headersMap {
		headersList = append(headersList, k)
	}
	for _, k := range headersList {
		v := headersMap[k]
		if len(k) > 6 && k[0:6] == "kafka_" {
			headersMap[k[6:]] = v
			delete(headersMap, k)
		}
	}
	content, err := args.Get("content")
	if err != nil {
		return nil, err
	}
	// convert content to []byte
	contentBytes, ok := content.([]byte)
	if !ok {
		return nil, err
	}
	// turn content which can be empty or not a json object into an encoded json string
	return func() (interface{}, error) {
		record := &kgo.Record{
			Value: contentBytes,
		}
		key, ok := headersMap["key"]
		if ok {
			record.Key = []byte(key.(string))
		}
		delete(headersMap, "key")

		// get timestamp
		record.Timestamp = time.Now()
		timestamp, ok := headersMap["kafka_timestamp_unix"]
		if ok {
			timestampUnix, ok := timestamp.(int64)
			if ok {
				record.Timestamp = time.Unix(timestampUnix, 0)
			}
		}
		delete(headersMap, "kafka_timestamp_unix")
		// put the remaining headers as kafka headers
		record.Headers = make([]kgo.RecordHeader, 0, len(headersMap))
		for k, v := range headersMap {
			sV, ok := v.(string)
			if !ok {
				sVA, ok := v.([]interface{})
				if ok {
					for _, v2 := range sVA {
						sV2, ok := v2.(string)
						if ok {
							record.Headers = append(record.Headers, kgo.RecordHeader{
								Key:   k,
								Value: []byte(sV2),
							})
						}
					}
				}
				continue
			}

			record.Headers = append(record.Headers, kgo.RecordHeader{
				Key:   k,
				Value: []byte(sV),
			})
		}

		kafkaClient := s.deadLetterClient.GetClient()
		var wg sync.WaitGroup
		wg.Add(1)
		var producerErr error
		kafkaClient.Produce(context.Background(), record, func(_ *kgo.Record, err error) {
			defer wg.Done()
			producerErr = err
			if producerErr != nil {
				log.Error().Err(producerErr).Msg("failed to produce")
			}
		})
		wg.Wait()
		if producerErr != nil {
			return nil, producerErr
		}
		return make(map[string]interface{}), nil

	}, nil
}
