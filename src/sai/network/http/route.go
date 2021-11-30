package http

import (
	"fmt"
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/saiset-co/saiStorageMongo/src/sai/auth"
)

type Validation struct {
	Key      string
	Type     string
	Rules    []string
	Required bool
}

//type Handler func(*fasthttp.RequestCtx)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(r *Route)
	Handlerr    func(*HttpRequest)
	Handler     fasthttp.RequestHandler
	RouteGroup  *RouteGroup
	Permission  auth.Rights
	Validations []Validation
	//RequestChannel  chan *http.Request
	//ResponseChannel chan *Response
}

func (route *Route) Register() {
	//route.RequestChannel = make(chan *http.Request)
	//route.ResponseChannel = make(chan *Response)
	//
	if route.RouteGroup == nil {
		route.RouteGroup = &RouteGroup{
			Chain: CreateMiddlewareChain(),
		}
	}

	API[route.Pattern] = route

	//go r.HandlerFunc(r)
}

func (route *Route) Handle(request *HttpRequest) {
	route.Handlerr(request)
}

func (r *Route) String() string {
	return fmt.Sprintf("%s [%s] %s", r.Name, r.Method, r.Pattern)
}
