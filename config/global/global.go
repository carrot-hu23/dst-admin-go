package global

import (
	"dst-admin-go/collect"
	"dst-admin-go/config"
	"net/http/httputil"
	"net/url"
)

var Collect *collect.Collect

type Route struct {
	Proxy *httputil.ReverseProxy
	Url   *url.URL
}

var RoutingTable = make(map[string]*Route)

var Config *config.Config
