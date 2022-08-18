package mockgen

import "github.com/splice/platform/services/license/rpc/licenseservice/v1"

func init() {
	registrar.AddService(&licenseservice.LicenseService{})
}
