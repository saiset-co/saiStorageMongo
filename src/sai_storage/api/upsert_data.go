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

func AddUpsertDataMethod() {
	route := &http.Route{
		Name:       "Upsert",
		Method:     "POST",
		Pattern:    "/upsert",
		Handler:    upsert,
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
			{
				Key:      "data",
				Type:     "object",
				Required: true,
			},
			{
				Key:      "options",
				Type:     "string",
				Required: false,
			},
		},
	}
	route.Register()
}

func upsert(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	collection := fmt.Sprint(request["collection"])
	dataArgument := request["data"]
	selectArgument := request["select"]
	options := request["options"]

	selector := map[string]interface{}{}
	// TODO add user id t oselect
	//if _, exist := request.GetParam("user_id"); exist {
	//	selector["user_id"] = session.Token.UserID
	//}
	if selectArgument != nil {
		storage.CopyToData(selector, selectArgument)
	}
	data := storage.CreateNewData(dataArgument)
	//TODO add user id to data
	//data["user_id"] = session.Token.UserID
	//data["user_id"] = request.Session.Token.UserID

	result := make([]interface{}, 0)

	if err := mongo.Upsert(fmt.Sprint(collection), selector, data, options, &result); err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		ctx.SetBody(common.ConvertInterfaceToJson(result))
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
}
