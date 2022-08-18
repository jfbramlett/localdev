package splicerouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"net/http"
	"testing"
	"time"

	"github.com/splice/platform/infra/libs/golang/internalhttp"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpliceRouter(t *testing.T) {
	router := NewSpliceRouter("infra/cmd/localdev/internal/splicerouter/testdata/test-routing.yaml", "infra/cmd/localdev/web")
	go func() {
		_ = router.Run(context.Background())
	}()

	go startTestRouter(8081)
	time.Sleep(2 * time.Second)

	client := internalhttp.NewGoHttpClient(context.Background(), 5*time.Second, internalhttp.WithInsecureSkipVerify(true))

	t.Run("test serve mock response matching query param", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/static-service/static-data")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		expected, err := os.ReadFile("infra/cmd/localdev/internal/splicerouter/testdata/static-service/get_static-data.mock")
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(body))
	})

	t.Run("test serve mock response matching path, param does not", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/static-service/static-param?arg=foo")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, `{"field": "default handler"}`, string(body))
	})

	t.Run("test serve mock response query param match", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/static-service/static-param?arg=mock")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		expected, err := os.ReadFile("infra/cmd/localdev/internal/splicerouter/testdata/static-service/get_static-param.mock")
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(body))
	})

	t.Run("test serve mock response header match", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "https://local.splice.com:8080/static-service/static-header", nil)
		require.NoError(t, err)

		req.Header.Add("x-query-loc", "mock")
		resp, err := client.Do(req)
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		expected, err := os.ReadFile("infra/cmd/localdev/internal/splicerouter/testdata/static-service/get_static-header.mock")
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(body))
	})

	t.Run("test serve mock response matching partial path", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/some-service/static-data")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		expected, err := os.ReadFile("infra/cmd/localdev/internal/splicerouter/testdata/some-service/get_static-data.mock")
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(body))
	})

	t.Run("test serve mock with direct config path", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/fixed-service/some/path/to/data")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		expected, err := os.ReadFile("infra/cmd/localdev/internal/splicerouter/testdata/direct_path_mock.mock")
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(body))
	})

	t.Run("test handlebar template", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/handlebar/assets?asset_uuids=12345&asset_uuids=54321")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		fmt.Println(string(body))
		var responseMap map[string]interface{}
		err = json.Unmarshal(body, &responseMap)
		require.NoError(t, err)

		responseAssets := responseMap["assets"]
		require.NotNil(t, responseAssets)
		assetsList, ok := responseAssets.([]interface{})
		require.True(t, ok)
		require.Len(t, assetsList, 2)
		asset1, ok := assetsList[0].(map[string]interface{})
		require.True(t, ok)
		asset2, ok := assetsList[1].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "12345", asset1["asset_uuid"])
		assert.True(t, asset1["asset_id"].(float64) > 0)
		assert.Equal(t, "54321", asset2["asset_uuid"])
		assert.True(t, asset2["asset_id"].(float64) > 0)
	})

	t.Run("test go template", func(t *testing.T) {
		resp, err := client.Get("https://local.splice.com:8080/catalog/assets?asset_uuids=12345&asset_uuids=54321")
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		require.NoError(t, err)

		fmt.Println(string(body))
		var responseMap map[string]interface{}
		err = json.Unmarshal(body, &responseMap)
		require.NoError(t, err)

		responseAssets := responseMap["assets"]
		require.NotNil(t, responseAssets)
		assetsList, ok := responseAssets.([]interface{})
		require.True(t, ok)
		require.Len(t, assetsList, 2)
		asset1, ok := assetsList[0].(map[string]interface{})
		require.True(t, ok)
		asset2, ok := assetsList[1].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "12345", asset1["asset_uuid"])
		assert.True(t, asset1["asset_id"].(float64) > 0)
		assert.Equal(t, "54321", asset2["asset_uuid"])
		assert.True(t, asset2["asset_id"].(float64) > 0)
	})

}

func startTestRouter(port int) {

	testRouter := mux.NewRouter()
	testRouter.PathPrefix("/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"field": "default handler"}`))
	})

	srv := &http.Server{
		Handler:      testRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	_ = srv.ListenAndServe()
}
