package mockgen

import "github.com/splice/platform/services/entitlement/rpc/entitlementservice/v1"

func init() {
	registrar.AddService(&entitlementservice.EntitlementService{})
}
