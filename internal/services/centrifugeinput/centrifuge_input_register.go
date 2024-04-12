package centrifugeinput

import (
	"context"
	"time"

	benthos_service "github.com/benthosdev/benthos/v4/public/service"
	centrifuge "github.com/centrifugal/centrifuge-go"
	contracts_storage "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/contracts/storage"
	utils "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/internal/utils"
	pkg_contracts_centrifuge "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/contracts/centrifuge"
	pkg_services_centrifuge_client "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/services/centrifuge/client"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	// Output
	InputName          = "centrifuge_input"
	FieldName_endpoint = "endpoint"
	FieldName_channel  = "channel"
)

var configSpec = benthos_service.NewConfigSpec().
	Summary("Creates an input of a centrifuge channel.").
	Field(benthos_service.NewStringField(FieldName_endpoint)).
	Field(benthos_service.NewStringField(FieldName_channel)).
	Field(benthos_service.NewBatchPolicyField("batching").
		Description("Allows you to configure a [batching policy](/docs/configuration/batching) that applies to individual topic partitions in order to batch messages together before flushing them for processing. Batching can be beneficial for performance as well as useful for windowed processing, and doing so this way preserves the ordering of topic partitions.").
		Advanced())

func (s *service) Register() error {
	ctx := context.Background()
	benthos_service.RegisterInput(InputName,
		configSpec, func(conf *benthos_service.ParsedConfig, mgr *benthos_service.Resources) (out benthos_service.Input, err error) {
			s.logger = mgr.Logger()
			log := s.logger.With("input", InputName).With("a", 1)
			log.With("caller", utils.Caller()).Info("Register")

			channel, err := conf.FieldString(FieldName_channel)
			if err != nil {
				return nil, err
			}
			s.channel = channel

			endpoint, err := conf.FieldString(FieldName_endpoint)
			if err != nil {
				return nil, err
			}
			s.endpoint = endpoint

			if s.batchPolicy, err = conf.FieldBatchPolicy("batching"); err != nil {
				return nil, err
			}
			var period time.Duration
			if s.batchPolicy.Period != "" {
				if period, err = time.ParseDuration(s.batchPolicy.Period); err != nil {

					return nil, status.Error(codes.InvalidArgument, "failed to parse duration string")
				}
			}
			s.period = period
			getLatestStreamPostitionResponse, err := s.centrifugeInputStorage.
				GetLatestStreamPostition(&contracts_storage.GetLatestStreamPostitionRequest{
					Namespace: s.channel,
				})
			if err != nil {
				s.log.Error().Err(err).Msg("failed to GetLatestStreamPostition")
				return nil, err
			}
			streamPosition := &centrifuge.StreamPosition{
				Epoch:  "0000", // this forces the stream to start at the known begining.
				Offset: 0,
			}
			if getLatestStreamPostitionResponse.StreamPosition != nil {
				streamPosition = getLatestStreamPostitionResponse.StreamPosition
			}
			s.centrifugeStreamBatcher.Configure(ctx, &pkg_contracts_centrifuge.CentrifugeConfig{
				Channel:                  "chat:index",
				BatchSize:                int32(s.batchPolicy.Count),
				NumberOfBatches:          2, // this is to cache ahead.
				HistoricalStreamPosition: streamPosition,
				CentrifugeClientConfig: &pkg_contracts_centrifuge.CentrifugeClientConfig{
					Endpoint: endpoint,
					GetToken: func(e centrifuge.ConnectionTokenEvent) (string, error) {
						return pkg_services_centrifuge_client.ExampleConnToken("49", 0), nil
					},
				},
			})

			return s, nil
		})
	return nil
}
