package api

import (
	"encoding/json"
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/common"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/db/mongo"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/network/http"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/storage"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai_storage/routing"
)

func AddRemoveDataMethod() {
	route := &http.Route{
		Name:       "Remove",
		Method:     "POST",
		Pattern:    "/remove",
		Handler:    remove,
		RouteGroup: routing.WithAuthAndValidation(),
		Validations: []http.Validation{
			{
				Key: "collection",
				Rules: []string{
					"string",
				},
				Type:     "string",
				Required: true,
			},
			{
				Key:      "select",
				Type:     "object",
				Required: false,
			},
		},
	}
	route.Register()
}

func remove(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	collection := fmt.Sprint(request["collection"])
	selectorArgument := request["select"]
	userId := request["user_id"]

	selector := map[string]interface{}{}
	if userId != nil {
		selector["user_id"] = "123"
	}
	if selectorArgument != nil {
		storage.CopyToData(selector, selectorArgument)
	}

	updatedData := []interface{}{}
	if err := mongo.Remove(fmt.Sprint(collection), selector, &updatedData); err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		ctx.SetBody(common.ConvertInterfaceToJson(updatedData))
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
}
