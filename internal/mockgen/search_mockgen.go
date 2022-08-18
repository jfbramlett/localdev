package mockgen

import "github.com/splice/platform/services/search/rpc/searchservice/v1"

func init() {
	registrar.AddService(&searchservice.SearchService{})
}
