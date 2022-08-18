package mockgen

import "github.com/splice/platform/services/taxonomy/rpc/taxonomyservice/v1"

func init() {
	registrar.AddService(&taxonomyservice.TaxonomyService{})
}
