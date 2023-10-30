package auth

import (
	proto_cloudeventprocessor "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/cloudeventprocessor"
	proto_kafkacloudevent "github.com/fluffy-bunny/benthos-cloudevents-fluffycore/pkg/proto/kafkacloudevent"

	contracts_common "github.com/fluffy-bunny/fluffycore/contracts/common"
	services_common_claimsprincipal "github.com/fluffy-bunny/fluffycore/services/common/claimsprincipal"
)

var writeEndpoints = []string{}
var noAuthEndpoints = []string{
	"/grpc.health.v1.Health/Check",
	// for no no auth
	proto_cloudeventprocessor.CloudEventProcessor_ProcessCloudEvents_FullMethodName,
	proto_kafkacloudevent.KafkaCloudEventService_SubmitCloudEvents_FullMethodName,
}

func BuildGrpcEntrypointPermissionsClaimsMap() map[string]contracts_common.IEntryPointConfig {
	entryPointClaimsBuilder := services_common_claimsprincipal.NewEntryPointClaimsBuilder()
	for _, endpoint := range noAuthEndpoints {
		entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(endpoint)
	}
	for _, endpoint := range writeEndpoints {
		entrypointConfig := &services_common_claimsprincipal.EntryPointConfig{
			FullMethodName: endpoint,
			ClaimsAST: &services_common_claimsprincipal.ClaimsAST{
				Or: []contracts_common.IClaimsValidator{
					&services_common_claimsprincipal.ClaimsAST{
						ClaimFacts: []contracts_common.IClaimFact{
							services_common_claimsprincipal.NewClaimFact(contracts_common.Claim{
								Type:  "permissions",
								Value: "write",
							}),
						},
					},
				},
			},
		}
		entryPointClaimsBuilder.EntrypointClaimsMap[endpoint] = entrypointConfig
	}
	entryPointClaimsBuilder.DumpExpressions()
	return entryPointClaimsBuilder.GetEntryPointClaimsMap()
}
