package splicerouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/splice/platform/localdev/v2/localdev/internal/files"
	"github.com/splice/platform/localdev/v2/localdev/internal/mockgen"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"

	"gopkg.in/yaml.v3"

	"github.com/aymerick/raymond"
	"github.com/sirupsen/logrus"
	"github.com/splice/platform/infra/libs/golang/internalhttp"
	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
)

const (
	RouteTypeProxy  = "proxy"
	RouteTypeStatic = "static"

	InternalCallPort = 8080
	ExternalCallPort = 8085

	SSLCert = `-----BEGIN CERTIFICATE-----
MIIEtDCCA5ygAwIBAgIJALgLRim5W7emMA0GCSqGSIb3DQEBCwUAMIGcMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCV0ExEDAOBgNVBAcMB1NlYXR0bGUxDzANBgNVBAoM
BlNwbGljZTEbMBkGA1UECwwSU3BsaWNlIEVuZ2luZWVyaW5nMRkwFwYDVQQDDBBs
b2NhbC5zcGxpY2UuY29tMSUwIwYJKoZIhvcNAQkBFhZlbmdpbmVlcmluZ0BzcGxp
Y2UuY29tMB4XDTIwMDcxNTE5MTQ1NVoXDTIyMTAxODE5MTQ1NVowgZwxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJXQTEQMA4GA1UEBwwHU2VhdHRsZTEPMA0GA1UECgwG
U3BsaWNlMRswGQYDVQQLDBJTcGxpY2UgRW5naW5lZXJpbmcxGTAXBgNVBAMMEGxv
Y2FsLnNwbGljZS5jb20xJTAjBgkqhkiG9w0BCQEWFmVuZ2luZWVyaW5nQHNwbGlj
ZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCzj2GOY0fjvbHI
AC1zMVEKXkXpTyb3/BIgSYTpoTBHPdcXOrkLLGmcvzHrsrX6HZ9Uf6QLfmCZ82qF
311m2a5M8bgOqnSJmsS4yWcFqiH1FmLMhM6vIqI0b/7zHUIPik6RmNcvlh09bFx0
j4tNc8l+0zbpGBwcrDQLKVv93nz2CPyzfgU7j/22yJ7Mre7nx90KNttFUeR0+aaR
+3u7HqmFgl5i/sUwW3lMd6JNRF9FYKxvYZCS8Ctenp4gmuHzOj5L7ovXJ4zfBPgU
b5nfsceXg2fOMNTNkRaafi8hHKxhOs9J2Mqrtx/EaSZEmmncBMq/TATMgmWzQRLY
muVKdZ0VAgMBAAGjgfYwgfMwgbsGA1UdIwSBszCBsKGBoqSBnzCBnDELMAkGA1UE
BhMCVVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdTZWF0dGxlMQ8wDQYDVQQKDAZT
cGxpY2UxGzAZBgNVBAsMElNwbGljZSBFbmdpbmVlcmluZzEZMBcGA1UEAwwQbG9j
YWwuc3BsaWNlLmNvbTElMCMGCSqGSIb3DQEJARYWZW5naW5lZXJpbmdAc3BsaWNl
LmNvbYIJAOJDTsadL8r3MAkGA1UdEwQCMAAwCwYDVR0PBAQDAgTwMBsGA1UdEQQU
MBKCEGxvY2FsLnNwbGljZS5jb20wDQYJKoZIhvcNAQELBQADggEBAHmi3NNfhxK1
qUJzEMxU5bU5oecC1MHKKMbQ99HzakwgDXBIJzcAYt02UE5b2oKFxvVpmU6dy2EU
jJ367rQCuH3IC3pnPDb34QlFeDodCdgG2sDgravrE+oT5R3Rm5l4kp62AcXKra8T
D1Q6BZazHbRHWKkq4XWHSqleAfJYg7xIOlKyitJb6NrswKH9RkOFtbQUC53rkYIR
owGbF0Ux/8xmZqrVrSFPMSWuAjmhSgiUAYqTD9EfQk84PiBhklj92J9qoPWS9w5s
p10JGRqItj2xPH1V7eOAEWh8IvT9O4pb60PEtPdy0zmiqc9VffSyJPmcKfxIP9N1
0HOCVkVGKJQ=
-----END CERTIFICATE-----`

	SSLPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAs49hjmNH472xyAAtczFRCl5F6U8m9/wSIEmE6aEwRz3XFzq5
CyxpnL8x67K1+h2fVH+kC35gmfNqhd9dZtmuTPG4Dqp0iZrEuMlnBaoh9RZizITO
ryKiNG/+8x1CD4pOkZjXL5YdPWxcdI+LTXPJftM26RgcHKw0Cylb/d589gj8s34F
O4/9tsiezK3u58fdCjbbRVHkdPmmkft7ux6phYJeYv7FMFt5THeiTURfRWCsb2GQ
kvArXp6eIJrh8zo+S+6L1yeM3wT4FG+Z37HHl4NnzjDUzZEWmn4vIRysYTrPSdjK
q7cfxGkmRJpp3ATKv0wEzIJls0ES2JrlSnWdFQIDAQABAoIBAATh4y8x9+RzZgAg
MRwuvRthENmVM2QO0JdcKGcJ4Pwu7EhPanSpUn3WnQ7hnx0b7RcpGVbOquXEvqLB
PhLr3FdvQyfy7pRHRw4XO9vlkmLNpwYUnmIYuPxgBzipFQDDK/u9gixjEox3A9SP
iqQif9oY+CdCZpFJBWlASKYQhKlT7YpsBdCJfsEorMq8OLVZog51CHsGOxa4XQ+1
Ev1DLm1slbPhQIMsrs0uLpB/dO4OIn5MfJK2SW/FGNmL6GvX3bbLmdca2FIHcF4q
nsNUKM60YolOehQCfGNLnJtAUHKn8Lk42C29dSFhz8Fq0SGjpsqQS+jTyw8c/IYD
ROXMr0ECgYEA5oJKfI+uDWTk4P8neY6yMeauEAdB0AKzoPf8k+iB373sddnJOhXL
nbOhEtHeni7JVSivgyxnfruOLGtvWWTEebP2Xxph45R2B2KEhqu1gcT+joBCwu8W
ebvF27O/K6Wmy6aYGMYTWn3AjxFngamm/+JYlbF7MughUP76pfr6MukCgYEAx2q6
QHlp8RZ/L/PhVKfCuA4dxuymVKs9uJcaY1IoLoqHYW4QbEyruyZYLWi4AmEkDDlW
j2IZv5d+r2N2hyfYqzb0WSuUCbRTaJApAIjmITV44/q/5moIxyMaByoAKvkCJV/L
ji+IWbLXjBr7YDrfmSBQdPpRDiyEqGwwEAlvxU0CgYEA5Ti9f56Vk6Y4YHH7PEsl
crAVecTtsj9th23zktYMiIViJlObYpKX98vQKlnfCeg2t+OMnWHDzWgPWqa/hOLK
6seGAU7H9zsEIBXc+dq41UIjbWuoeBavgAC1IeRd/7Zr5mpVJ5WZW0xf9yV0i6E8
e4sHUly6yYXC07urXvD0azkCgYAJ0nvyCQjq0wzYs855ePniTu+wiJ94tCaKHQcz
tSw9fp1Ec0Nj0jLzOORG+E138IjyATD+Rvq1sSSQRvnjllbZuA85BSh5geRJ1i/u
0s9i+1tE/2jMVJSyGkyB5dO0SieM57cC/dxdbq2nPPz8tGmnBSxxVpL/e7ndAdcs
MwrKUQKBgG2AARuK2SQZPYqe6DL/0lwmFFQIHFuQ3Ip4kcipLgs6gUrd/rQScv+b
CCblMWJ54s5Tyrl0bWB4Gn1ZBq9XO51xxEc4rEXStpDYmFm6Uuysnug8bHHgTk6m
8YAQEhAgeI53zA8aSXq7WQm/WuImNyQ9IZAjiKGUfzKm/019mXPz
-----END RSA PRIVATE KEY-----`
)

const (
	templateTypeGo        = "go"
	templateTypeHandlebar = "handlebar"
)

func loadConfig(filename string) (*Configuration, error) {
	if filename == "" {
		filename = RouterConfig
	}
	var srCfg Configuration
	if fi, err := os.Stat(filename); err == nil {
		if !fi.IsDir() {
			file, err := files.ReadAll(filename)
			if err != nil {
				return nil, err
			}

			var srCfg Configuration
			if err := yaml.Unmarshal([]byte(file), &srCfg); err != nil {
				return nil, err
			}
		}
	}

	for _, route := range srCfg.Routes {
		route.ID = uuid.NewString()
	}
	srCfg.SortRoute()

	return &srCfg, nil
}

// NewSpliceRouter creates a new instance of our splice router
func NewSpliceRouter(configFile string, webRoot string) *SpliceRouter {
	var err error
	var sslCrtFile string
	var sslKeyFile string
	if sslCrtFile, err = files.CreateFileInTmpDir("local.splice.com.crt", SSLCert); err != nil {
		os.Exit(1)
	}
	if sslKeyFile, err = files.CreateFileInTmpDir("local.splice.com.key", SSLPrivateKey); err != nil {
		os.Exit(1)
	}

	cfg, err := loadConfig(configFile)
	if err != nil {
		os.Exit(1)
	}

	return &SpliceRouter{
		config:         cfg,
		sslCertFile:    sslCrtFile,
		sslCertKeyFile: sslKeyFile,
		webRoot:        webRoot,
		mockGenerator:  mockgen.NewMockGenerator(),
	}
}

// SpliceRouter provides an implementation of a dynamic router. Much like an API Gateway this just routes requests - it differs
// from something like Kong in it is not designed to enrich a request it is pure pass-thru. It can also serve static (i.e. mock)
// content. The goal is something simple, with a basic config that let's us run this to route traffic between mock, local services,
// and services in a shared environment.
// The configuration supports "internal" and an "external" routing definitions. This is done to support service-to-service
// routing which don't go through another gateway where as external calls (like from GraphQL) do since the routes might be
// the same.
type SpliceRouter struct {
	sslCertFile    string
	sslCertKeyFile string
	config         *Configuration
	webRoot        string
	mockGenerator  *mockgen.MockGen
}

func (sr *SpliceRouter) Run(ctx context.Context) error {
	go sr.RunAdmin(ctx)

	raymond.RegisterHelper("is_last", IsLast)

	var wait sync.WaitGroup

	if sr.config != nil {
		wait.Add(1)
		go func() {
			internalCtx, log := splicelogger.UpdateInContext(ctx, logrus.Fields{"config": "internal"})
			err := sr.configureListener(internalCtx, InternalCallPort, sr.config)
			if err != nil {
				log.WithError(err).Error("failed configuring our internal listener")
			}
			wait.Done()
		}()
	}

	wait.Wait()

	return nil
}

func (sr *SpliceRouter) configureListener(ctx context.Context, listenPort int, cfg *Configuration) error {
	logger, _ := splicelogger.FromContext(ctx)

	for _, route := range cfg.Routes {
		_ = sr.mapRouteHandler(ctx, route)
	}

	srv := &http.Server{
		Handler: sr,
		Addr:    fmt.Sprintf(":%d", listenPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	logger.Infof("starting listener on port %d with cert %s and cert key %s", listenPort, sr.sslCertFile, sr.sslCertKeyFile)
	if err := srv.ListenAndServeTLS(sr.sslCertFile, sr.sslCertKeyFile); err != nil {
		return err
	}

	return nil
}

func (sr *SpliceRouter) mapRouteHandler(ctx context.Context, route *Route) error {
	logger, _ := splicelogger.FromContext(ctx)
	route.regex = regexp.MustCompile("^" + route.Path + "$")
	switch route.Type {
	case RouteTypeProxy:
		proxy, err := newReverseProxy(route.Destination)
		if err != nil {
			return err
		}

		route.handler = func(writer http.ResponseWriter, request *http.Request) {
			logger.WithField("path", request.URL.Path).Info("routing to proxy")
			for k, v := range route.HeadersEnrichment {
				request.Header.Add(k, v)
			}
			proxy.ServeHTTP(writer, request)
		}
	case RouteTypeStatic:
		proxy, err := newStaticFileHandler(route)
		if err != nil {
			return err
		}
		route.handler = func(writer http.ResponseWriter, request *http.Request) {
			logger.WithField("path", request.URL.Path).Info("routing to mock")
			proxy.ServeHTTP(writer, request)
		}
	}

	return nil
}

func (sr *SpliceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range sr.config.Routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if route.Method != "" {
				if r.Method != route.Method {
					allow = append(allow, route.Method)
					continue
				}
			}

			matchHeaders := true
			for k, v := range route.HeadersMatch {
				if v != r.Header.Get(k) {
					matchHeaders = false
					break
				}
			}
			if !matchHeaders {
				continue
			}

			matchParams := true
			queryParams := r.URL.Query()
			for k, v := range route.QueryParams {
				if v != queryParams.Get(k) {
					matchParams = false
					break
				}
			}
			if !matchParams {
				continue
			}

			route.handler(w, r)
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func newStaticFileHandler(route *Route) (http.HandlerFunc, error) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error
		mockContent := route.MockData
		if mockContent == "" {
			fileName := route.Destination
			if !files.Exists(fileName) {
				pathElements := strings.Split(request.URL.Path, "/")
				name := pathElements[len(pathElements)-1]
				fileName = fmt.Sprintf("%s/%s_%s.mock", route.Destination, strings.ToLower(request.Method), name)
				if !files.Exists(fileName) {
					writeResponse(writer, http.StatusBadGateway, fmt.Sprintf(`{"status": "failed to find mock file", "file": "%s"}`, fileName))
					return
				}
			}

			mockContent, err = files.ReadAll(fileName)
			if err != nil {
				writeResponse(writer, http.StatusBadGateway, fmt.Sprintf(`{"status": "failed to find mock file", "file": "%s", "error": "%v"}`, fileName, err))
				return
			}
		}

		switch route.TemplateType {
		case templateTypeGo:
			templateData, err := getTemplateArgs(request)
			if err != nil {
				writeResponse(writer, http.StatusBadGateway, fmt.Sprintf(`{"status": "failed to find mock file", "file": "%s", "error": "%v"}`, route.Destination, err))
				return
			}
			t, err := template.New(route.Destination).Funcs(map[string]interface{}{
				"is_last":    IsLast,
				"random_int": RandomInt,
			}).Parse(mockContent)
			if err != nil {
				writeResponse(writer, http.StatusOK, mockContent)
				return
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, templateData); err != nil {
				writeResponse(writer, http.StatusOK, mockContent)
				return
			}

			writeResponse(writer, http.StatusOK, buf.String())
			return
		case templateTypeHandlebar:
			templateData, err := getTemplateArgs(request)
			if err != nil {
				writeResponse(writer, http.StatusBadGateway, fmt.Sprintf(`{"status": "failed to generate template data", "error": "%v"}`, err))
				return
			}

			templateData["random_int"] = RandomInt

			resp, err := raymond.Render(mockContent, templateData)
			if err != nil {
				writeResponse(writer, http.StatusBadGateway, fmt.Sprintf(`{"status": "failed to render mock file", "file": "%s", "error": "%v"}`, route.Destination, err))
				return
			}

			writeResponse(writer, http.StatusOK, resp)
			return
		default:
			writeResponse(writer, http.StatusOK, mockContent)
			return
		}

	}, nil
}

func getTemplateArgs(request *http.Request) (map[string]interface{}, error) {
	templateData := make(map[string]interface{})
	switch request.Method {
	case http.MethodPost:
		bodyStr, err := io.ReadAll(request.Body)
		defer func() {
			_ = request.Body.Close()
		}()

		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(bodyStr, &templateData)
	case http.MethodGet:
		args := request.URL.Query()
		for k, vals := range args {
			if len(vals) > 1 {
				templateData[k] = vals
			} else {
				templateData[k] = vals[0]
			}
		}
	}

	return templateData, nil
}

func newReverseProxy(location string) (http.Handler, error) {
	u, err := url.Parse(location)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		ErrorHandler: func(w http.ResponseWriter, req *http.Request, err2 error) {
			writeResponse(w, http.StatusBadGateway, `{"status": "service unavailable - Bad Gateway"}`)
		},
		Transport: internalhttp.NewInternalTransport(context.Background(), internalhttp.NewDefaultTransportConfig(internalhttp.WithInsecureSkipVerify(true))),
		Director: func(req *http.Request) {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.Host = u.Host
		},
	}, nil
}

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write([]byte(msg))
}
