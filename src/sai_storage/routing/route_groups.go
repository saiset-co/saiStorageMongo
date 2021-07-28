package routing

import (
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/network/http"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai_storage/middleware"
)

func WithAuth() *http.RouteGroup {
	return &http.RouteGroup{
		Chain: http.CreateMiddlewareChain(middleware.CheckAuth, middleware.CheckPermissions),
	}
}

func WithAuthAndValidation() *http.RouteGroup {
	return &http.RouteGroup{
		Chain: http.CreateMiddlewareChain(middleware.CheckAuth, middleware.CheckPermissions, middleware.ValidateRequest),
		//Chain: http.CreateMiddlewareChain(middleware.ValidateRequest, middleware.CheckAuth, middleware.CheckPermissions, middleware.EnableCorsHeader),
	}
}

func WithValidation() *http.RouteGroup {
	return &http.RouteGroup{
		Chain: http.CreateMiddlewareChain(middleware.ValidateRequest),
		//Chain: http.CreateMiddlewareChain(middleware.ValidateRequest, middleware.EnableCorsHeader),
	}
}
