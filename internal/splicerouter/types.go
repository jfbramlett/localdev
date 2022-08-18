package splicerouter

import (
	"net/http"
	"regexp"
	"sort"
)

type SimpleResult struct {
	Page  string
	Msg   string
	Error error
}

type Routes struct {
	Page   string
	Type   string
	Msg    string
	Error  error
	Routes []*Route
}

type Route struct {
	Priority          int               `yaml:"priority"`
	Method            string            `yaml:"method"`
	HeadersMatch      map[string]string `yaml:"headers_match"`
	QueryParams       map[string]string `yaml:"query_params"`
	Path              string            `yaml:"path"`
	Destination       string            `yaml:"destination"`
	Type              string            `yaml:"type"`
	TemplateType      string            `yaml:"template_type"`
	MockData          string            `yaml:"mock_data"`
	HeadersEnrichment map[string]string `yaml:"headers_enrichment"`
	regex             *regexp.Regexp    `yaml:"-"`
	handler           http.HandlerFunc  `yaml:"-"`
	ID                string            `yaml:"-"`
}

type Configuration struct {
	Routes []*Route `yaml:"routes"`
}

func (c *Configuration) SortRoute() {
	sort.SliceStable(c.Routes, func(i, j int) bool {
		return c.Routes[i].Priority < c.Routes[j].Priority
	})
}

type Mocks struct {
	MockTypes []string
	Type      string
	MockData  string
	Error     error
	Msg       string
}
