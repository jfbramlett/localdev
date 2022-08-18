package mockgen

import "github.com/splice/platform/services/merchandising/rpc/merchandisingservice/v1"

func init() {
	registrar.AddService(&merchandisingservice.MerchandisingService{})
}
