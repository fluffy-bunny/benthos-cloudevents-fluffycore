package bloblang

import (
	"encoding/json"

	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/rs/zerolog/log"
)

type (
	wrappedContent struct {
		Type    string      `json:"type"`
		Headers interface{} `json:"headers"`
		Content interface{} `json:"content"`
	}
)

func ddd() {
	crazyObjectSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewAnyParam("headers")).Param(bloblang.NewAnyParam("content"))

	err := bloblang.RegisterFunctionV2("wrap", crazyObjectSpec, func(args *bloblang.ParsedParams) (bloblang.Function, error) {
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
				Type:    "error",
				Headers: headers,
				Content: content,
			}
			obj2 := map[string]interface{}{}
			bb, _ := json.Marshal(obj)
			json.Unmarshal(bb, &obj2)
			log.Info().Interface("obj", obj).Msg("wrap")

			return obj2, nil

		}, nil
	})
	if err != nil {
		panic(err)
	}

	intoObjectSpec := bloblang.NewPluginSpec().
		Param(bloblang.NewStringParam("key"))

	err = bloblang.RegisterMethodV2("into_object", intoObjectSpec, func(args *bloblang.ParsedParams) (bloblang.Method, error) {
		key, err := args.GetString("key")
		if err != nil {
			return nil, err
		}

		return func(v interface{}) (interface{}, error) {
			return map[string]interface{}{key: v}, nil
		}, nil
	})
	if err != nil {
		panic(err)
	}
}
