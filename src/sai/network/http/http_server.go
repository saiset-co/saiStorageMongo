package http

import (
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/kirillbeldyaga/fasthttprouter"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/fatih/color"
	"log"
	"net/http"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/common"
)

type HttpServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (httpServer *HttpServer) Address() string {
	return fmt.Sprintf("%s:%d", httpServer.Host, httpServer.Port)
}

var (
	SaiHttpServer = &HttpServer{}
)

func SetHttpServerInv(httpServer HttpServer) {
	*SaiHttpServer = httpServer
}

func NotFound(ctx *fasthttp.RequestCtx) {
	err := RouteNotFindError(string(ctx.Path()))
	ctx.Error(err.Error(), err.Code)
	ctx.Response.Header.Set("Content-Type", "application/json")
}

func MethodNotAllowed(ctx *fasthttp.RequestCtx) {
	err := MethodNotAllowedError(string(ctx.Method()))
	ctx.Error(err.Error(), err.Code)
	ctx.Response.Header.Set("Content-Type", "application/json")
}

func SetResponse(ctx *fasthttp.RequestCtx, resp []byte) {
	ctx.SetStatusCode(200)
	ctx.SetBody(resp)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
}

func SetErrorResponse(ctx *fasthttp.RequestCtx, err *common.Error) {
	ctx.SetStatusCode(err.Code)
	ctx.SetBodyString(err.Error())

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
}

func (httpServer *HttpServer) Start() {
	router := fasthttprouter.New()
	router.NotFound = NotFound
	router.MethodNotAllowed = MethodNotAllowed

	RegisterHandlers(router)

	//mux.Handle("/", NoteTime(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	//	route, err := GetRoute(req.URL.Path)
	//	if err != nil {
	//		var response Response
	//		response.SetError(err)
	//
	//		res.WriteHeader(response.Code)
	//		res.Write(response.Body)
	//
	//		return
	//	}
	//
	//	request, err := GetRequest(req, route);
	//	if err != nil {
	//		var response Response
	//		response.SetError(err)
	//
	//		res.WriteHeader(response.Code)
	//		res.Write(response.Body)
	//
	//		return
	//	}
	//
	//	//go route.Handle(request, res)
	//	route.Handle(request)
	//	response := request.Response
	//
	//	if response.Error != nil {
	//		res.WriteHeader(response.Code)
	//	}
	//	res.Write(response.Body)
	//
	//	//response := <-request.ResponseChannel
	//
	//	//handler := route.RouteGroup.Chain.ThenFunc(func(w http.ResponseWriter, r *HttpRequest) {
	//	//	err := checkHttpMethod(w, r)
	//	//	if err != nil {
	//	//		return
	//	//	} else {
	//	//		route.RequestChannel <- r
	//	//		response := <-route.ResponseChannel
	//	//
	//	//		if response.Error != nil {
	//	//			w.WriteHeader(response.Code)
	//	//		}
	//	//		w.Write(response.Body)
	//	//	}
	//	//})
	//	//handler.ServeHTTP(res, req)
	//})))

	d := color.New(color.FgCyan, color.Bold)
	d.Println("********************************************")
	d.Println("Listening at", httpServer.Address())

	log.Fatal(fasthttp.ListenAndServe(httpServer.Address(), router.Handler))
}

func checkHttpMethod(res http.ResponseWriter, req *http.Request) error {
	if req.Method != "POST" && req.Method != "GET" {
		//TODO print to web
		http.Error(res, "Invalid request method", http.StatusMethodNotAllowed)
		// TODO return error
	}

	return nil
}
