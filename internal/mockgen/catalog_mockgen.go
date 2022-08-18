package mockgen

import "github.com/splice/platform/services/catalog/rpc/catalogservice/v1"

func init() {
	registrar.AddService(&catalogservice.CatalogService{})
}
