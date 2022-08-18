package mockgen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockGen(t *testing.T) {
	gen := NewMockGenerator()
	_, err := gen.GenerateMock("splice.api.taxonomy.v1.GetTaxonomyResponse")
	require.NoError(t, err)
}
