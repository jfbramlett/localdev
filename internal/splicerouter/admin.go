package splicerouter

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	splicelogger "github.com/splice/platform/infra/libs/golang/logger"
	"gopkg.in/yaml.v3"
)

const (
	PORT         = "8001"
	RouterConfig = "./router.yaml"
)

var (
	homeTpl         *template.Template
	routerTpl       *template.Template
	createRoutesTpl *template.Template
	simpleTpl       *template.Template
	genMocksTpl     *template.Template
)

func (sr *SpliceRouter) RunAdmin(ctx context.Context) {
	base := template.Must(template.ParseFiles(path.Join(sr.webRoot, "index.html"), path.Join(sr.webRoot, "templates", "shared", "router_nav_bar.html")))
	homeTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(sr.webRoot, "templates", "home.html")))
	routerTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(sr.webRoot, "templates", "router.html")))
	createRoutesTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(sr.webRoot, "templates", "create_route.html")))
	genMocksTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(sr.webRoot, "templates", "mock_gen.html")))
	simpleTpl = template.Must(template.Must(base.Clone()).ParseFiles(path.Join(sr.webRoot, "templates", "simple.html")))

	// set up HTTP server
	mux := http.NewServeMux()
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.HandleFunc("/", sr.indexHandler)
	mux.HandleFunc("/router", sr.routerHandler)
	mux.HandleFunc("/create/", sr.routesHandler)
	mux.HandleFunc("/create-new", sr.createRouteHandler)
	mux.HandleFunc("/mocks/", sr.mocksHandler)
	mux.HandleFunc("/generate-mock", sr.generateMocksHandler)
	mux.HandleFunc("/export-routes", sr.exportRoutesHandler)
	mux.HandleFunc("/update-route", sr.updateRouteHandler)
	mux.HandleFunc("/delete-route", sr.deleteRouteHandler)

	logger, _ := splicelogger.FromContext(ctx)

	logger.Infof("starting admin on port %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, mux))
}

func (sr *SpliceRouter) indexHandler(w http.ResponseWriter, r *http.Request) {
	result := &SimpleResult{}

	if r.Method == "POST" {
		result.Msg = "POST not supported"
		_ = simpleTpl.Execute(w, result)
	} else if r.Method == "GET" {
		err := homeTpl.Execute(w, nil)
		if err != nil {
			result := &SimpleResult{}
			result.Msg = err.Error()
			_ = simpleTpl.Execute(w, result)
		}

		return
	}
	_ = simpleTpl.Execute(w, result)
}

func (sr *SpliceRouter) routerHandler(w http.ResponseWriter, r *http.Request) {
	result := &SimpleResult{}

	if r.Method == "POST" {
		result.Msg = "POST not supported"
		_ = simpleTpl.Execute(w, result)
	} else if r.Method == "GET" {
		routes := Routes{Routes: sr.config.Routes, Type: "Routes", Page: "routes"}

		err := routerTpl.Execute(w, routes)
		if err != nil {
			result := &SimpleResult{}
			result.Msg = err.Error()
			_ = simpleTpl.Execute(w, result)
		}

		return
	}
	_ = simpleTpl.Execute(w, result)
}

func (sr *SpliceRouter) routesHandler(w http.ResponseWriter, r *http.Request) {
	err := createRoutesTpl.Execute(w, SimpleResult{Page: "create"})
	if err != nil {
		result := &SimpleResult{}
		result.Msg = err.Error()
		_ = simpleTpl.Execute(w, result)
	}
}

func (sr *SpliceRouter) createRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := &SimpleResult{
			Msg: "GET not supported",
		}
		_ = simpleTpl.Execute(w, result)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			result := &SimpleResult{
				Error: err,
			}
			_ = simpleTpl.Execute(w, result)
			return
		}

		newRoute := &Route{ID: uuid.NewString()}
		priority := 0
		priority, _ = strconv.Atoi(r.Form.Get("priority"))
		newRoute.Priority = priority
		newRoute.Path = r.Form.Get("path")
		newRoute.Method = r.Form.Get("method")
		newRoute.Type = r.Form.Get("proxy")

		headerEnrichments := r.Form.Get("headerEnrichments")
		if len(headerEnrichments) > 0 {
			headerSet := strings.Split(headerEnrichments, ";")
			for _, h := range headerSet {
				newRoute.HeadersEnrichment = make(map[string]string)
				kv := strings.Split(h, "=")
				if len(kv) == 2 {
					newRoute.HeadersEnrichment[kv[0]] = kv[1]
				}
			}
		}

		headers := r.Form.Get("headers")
		if len(headers) > 0 {
			headerSet := strings.Split(headers, ";")
			for _, h := range headerSet {
				newRoute.HeadersMatch = make(map[string]string)
				kv := strings.Split(h, "=")
				if len(kv) == 2 {
					newRoute.HeadersMatch[kv[0]] = kv[1]
				}
			}
		}

		proxyType := r.Form.Get("proxyType")
		switch proxyType {
		case RouteTypeProxy:
			newRoute.Type = RouteTypeProxy
			newRoute.Destination = r.Form.Get("reverseProxyDestination")
		case RouteTypeStatic:
			newRoute.Type = RouteTypeStatic
			newRoute.Destination = r.Form.Get("mockDataDestination")
			newRoute.TemplateType = r.Form.Get("templateType")
			newRoute.MockData = r.Form.Get("mockResponse")
		}

		_ = sr.mapRouteHandler(r.Context(), newRoute)
		sr.config.Routes = append(sr.config.Routes, newRoute)
		sr.config.SortRoute()

		http.Redirect(w, r, "/router", http.StatusSeeOther)

		return
	}

	result := &SimpleResult{
		Msg: "Invalid http verb",
	}
	_ = simpleTpl.Execute(w, result)
}

func (sr *SpliceRouter) mocksHandler(w http.ResponseWriter, r *http.Request) {
	mocks := &Mocks{
		MockTypes: sr.mockGenerator.ListTypes(),
	}
	if err := genMocksTpl.Execute(w, mocks); err != nil {
		result := &SimpleResult{}
		result.Msg = err.Error()
		_ = simpleTpl.Execute(w, result)
		return
	}
}

func (sr *SpliceRouter) generateMocksHandler(w http.ResponseWriter, r *http.Request) {
	mocks := &Mocks{
		MockTypes: sr.mockGenerator.ListTypes(),
	}

	if r.Method == "GET" {
		mocks.Msg = "GET not supported"
		_ = genMocksTpl.Execute(w, mocks)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			mocks.Error = err
			_ = genMocksTpl.Execute(w, mocks)
			return
		}

		typ := r.Form.Get("mocks")
		if typ == "" {
			mocks.Msg = "no type specified"
			_ = genMocksTpl.Execute(w, mocks)
			return
		}

		output, err := sr.mockGenerator.GenerateMock(typ)
		if err != nil {
			mocks.Error = err
			_ = genMocksTpl.Execute(w, mocks)
			return
		}

		mocks.Type = typ
		mocks.MockData = output
		_ = genMocksTpl.Execute(w, mocks)
		return
	}

	mocks.Msg = "Invalid http verb"
	_ = genMocksTpl.Execute(w, mocks)
}

func (sr *SpliceRouter) exportRoutesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := &SimpleResult{
			Msg: "GET not supported",
		}
		_ = simpleTpl.Execute(w, result)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			result := &SimpleResult{
				Error: err,
			}
			_ = simpleTpl.Execute(w, result)
			return
		}

		routesExport, _ := yaml.Marshal(sr.config)
		if _, err := os.Stat(RouterConfig); err == nil {
			_ = os.Rename(RouterConfig, fmt.Sprintf("./router-%d.yaml", time.Now().Nanosecond()))
		}
		_ = os.WriteFile(RouterConfig, routesExport, 0644)

		http.Redirect(w, r, "/router", http.StatusSeeOther)
		return
	}

	result := &SimpleResult{
		Msg: "Invalid http verb",
	}
	_ = simpleTpl.Execute(w, result)
}

func (sr *SpliceRouter) deleteRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := &SimpleResult{
			Msg: "GET not supported",
		}
		_ = simpleTpl.Execute(w, result)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			result := &SimpleResult{
				Error: err,
			}
			_ = simpleTpl.Execute(w, result)
			return
		}

		routeID := r.Form.Get("delete_id")
		newRoutes := make([]*Route, 0)
		for _, route := range sr.config.Routes {
			if route.ID != routeID {
				newRoutes = append(newRoutes, route)
			}
		}

		sr.config.Routes = newRoutes

		http.Redirect(w, r, "/router", http.StatusSeeOther)
		return
	}

	result := &SimpleResult{
		Msg: "Invalid http verb",
	}
	_ = simpleTpl.Execute(w, result)
}

func (sr *SpliceRouter) updateRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := &SimpleResult{
			Msg: "GET not supported",
		}
		_ = simpleTpl.Execute(w, result)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			result := &SimpleResult{
				Error: err,
			}
			_ = simpleTpl.Execute(w, result)
			return
		}

		routeID := r.Form.Get("update_id")
		for _, route := range sr.config.Routes {
			if route.ID == routeID {
				route.Method = r.Form.Get("update_method")
				route.Path = r.Form.Get("update_path")
				route.Destination = r.Form.Get("update_destination")
				priority, _ := strconv.Atoi(r.Form.Get("update_priority"))
				route.Priority = priority
				_ = sr.mapRouteHandler(r.Context(), route)
			}
		}

		sr.config.SortRoute()

		http.Redirect(w, r, "/router", http.StatusSeeOther)
		return
	}

	result := &SimpleResult{
		Msg: "Invalid http verb",
	}
	_ = simpleTpl.Execute(w, result)
}
