package api

import (
	"saiStorageMongo/src/sai/network/http"
	"saiStorageMongo/src/sai_storage/routing"
	"fmt"
	"saiStorageMongo/src/sai/db/mongo"
	"saiStorageMongo/src/sai/storage"
	"saiStorageMongo/src/sai/common"
	"saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"encoding/json"
)

func AddGetDataMethod() {
	route := http.Route{
		Name:       "Get",
		Method:     "POST",
		Pattern:    "/get",
		Handler:    get,
		Handlerr:   gett,
		RouteGroup: routing.WithAuthAndValidation(),
		//Permission: auth.Permission{
		//	Read:  7,
		//	Write: 7,
		//},
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

func get(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	collection := fmt.Sprint(request["collection"])
	selectorArgument := request["select"]
	userId := request["user_id"]
	options := request["options"]

	selector := map[string]interface{}{}

	if userId != nil {
		selector["user_id"] = "123"
	}
	if selectorArgument != nil {
		storage.CopyToData(selector, selectorArgument)
	}

	result := make([]interface{}, 0)
	err := mongo.Find(fmt.Sprint(collection), selector, options, &result)

	if err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		ctx.SetBody(common.ConvertInterfaceToJson(result))
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
}

func gett(request *http.HttpRequest) {
	var response http.Response

	collection, _ := request.GetParam("collection")
	selectorArgument, _ := request.GetParam("select")
	options,_ := request.GetParam("options")

	selector := map[string]interface{}{}
	if _, exist := request.GetParam("user_id"); exist {
		selector["user_id"] = request.Session.Token.UserID
	}
	if selectorArgument != nil {
		storage.CopyToData(selector, selectorArgument)
	}

	result := make([]interface{}, 0)
	err := mongo.Find(fmt.Sprint(collection), selector, options, &result)

	if err != nil {
		response.SetError(err)
	} else {
		response.SetBody(common.ConvertInterfaceToJson(result))
		response.Code = 200
	}
	request.ResponseChannel <- &response
}
