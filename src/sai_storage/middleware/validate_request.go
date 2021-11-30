package middleware

import (
	//http2 "sai/network/http"
	"encoding/json"
	"reflect"
	//"fmt"
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/saiset-co/saiStorageMongo/src/sai/network/http"
)

func ValidateRequest(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var request map[string]interface{}
		json.Unmarshal(ctx.PostBody(), &request)

		url := ctx.Request.URI().Path()
		route, _ := http.API[string(url)]
		for _, validation := range route.Validations {
			if param, exist := request[validation.Key]; exist {
				if !Satisfy(param, validation.Rules, validation.Type) {
					err := http.BadRequestError()
					http.SetErrorResponse(ctx, err)
					return
				}
			} else {
				if validation.Required {
					err := http.BadRequestError()
					http.SetErrorResponse(ctx, err)
					return
				}
			}
		}
		h(ctx)
	}
}

func Satisfy(param interface{}, rules []string, paramType string) bool {
	switch paramType {
	case "object":
		if _, err := json.Marshal(param); err != nil {
			return false
		} else {
			if reflect.TypeOf(param).Kind() != reflect.Map {
				return false
			}
		}
	case "string":
		if reflect.TypeOf(param).Kind() != reflect.String {
			return false
		}
	case "int":
		if reflect.TypeOf(param).Kind() != reflect.Int {
			return false
		}
	case "float":
		if reflect.TypeOf(param).Kind() != reflect.Float64 {
			return false
		}
	}
	for _, rule := range rules {
		if rule == "!empty" {
			if param == "" {
				return false
			}
		}
	}
	return true
}
