package bloblang

import (
	"github.com/benthosdev/benthos/v4/public/bloblang"
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
type (
	wrappedContent struct {
		Error   string      `json:"error"`
		Headers interface{} `json:"headers"`
		Content interface{} `json:"content"`
	}
)

func (s *service) WrapIt(args *bloblang.ParsedParams) (bloblang.Function, error) {
	headers, err := args.Get("headers")
	if err != nil {
		return nil, err
	}
	content, err := args.Get("content")
	if err != nil {
		return nil, err
	}

	// turn content which can be empty or not a json object into an encoded json string
	return func() (interface{}, error) {
		obj := &wrappedContent{
			Error:   "error",
			Headers: headers,
			Content: content,
		}
		return obj, nil

	}, nil
}
