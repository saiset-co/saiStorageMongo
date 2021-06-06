package http

import (
	"saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
)

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler

type MiddlewareChain struct {
	middlewares []Middleware
}

//func (middleware *Middleware) Register() {
//	Middlewares = append(Middlewares, middleware)
//}

func CreateMiddlewareChain(middlewares ...Middleware) *MiddlewareChain {
	return &MiddlewareChain{append(([]Middleware)(nil), middlewares...)}
}

func (c *MiddlewareChain) Then(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	if h == nil {
		h = func(ctx *fasthttp.RequestCtx) {}
	}

	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}

	return h
}

//func (c *MiddlewareChain) ThenFunc(fn Handler) http.Handler {
//	if fn == nil {
//		return c.Then(nil)
//	}
//	return c.Then(fn)
//}

//func (c *MiddlewareChain) ThenF(fn http.HandlerFunc) http.Handler {
//	if fn == nil {
//		return c.Then(nil)
//	}
//	return c.Then(fn)
//}

func (c *MiddlewareChain) Append(middlewares ...Middleware) *MiddlewareChain {
	newCons := make([]Middleware, 0, len(c.middlewares)+len(middlewares))
	newCons = append(newCons, c.middlewares...)
	newCons = append(newCons, middlewares...)

	return &MiddlewareChain{newCons}
}

func (c *MiddlewareChain) Extend(chain MiddlewareChain) *MiddlewareChain {
	return c.Append(chain.middlewares...)
}
