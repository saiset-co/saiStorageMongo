package http

import (
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttprouter"
)

var (
	API = map[string]*Route{}
)

type RouteGroup struct {
	Chain *MiddlewareChain
}

func RegisterHandlers(router *fasthttprouter.Router) {
	for _, route := range API {
		if route.Method == "GET" {
			router.GET(route.Pattern, route.RouteGroup.Chain.Then(route.Handler))
		} else {
			router.POST(route.Pattern, route.RouteGroup.Chain.Then(route.Handler))
		}
	}
}
