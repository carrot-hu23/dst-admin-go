package entity

import (
	"net/http/httputil"
	"net/url"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Route struct {
	Proxy *httputil.ReverseProxy
	Url   *url.URL
}

var RoutingTable = make(map[string]*Route)
