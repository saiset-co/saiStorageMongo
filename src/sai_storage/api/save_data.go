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

func AddSaveDataMethod() {
	route := &http.Route{
		Name:       "Save",
		Method:     "POST",
		Pattern:    "/save",
		Handler:    save,
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
				Key:      "data",
				Type:     "object",
				Required: true,
			},
		},
	}
	route.Register()
}

func save(ctx *fasthttp.RequestCtx) {
	//var response http.Response

	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	collection := fmt.Sprint(request["collection"])
	dataArgument := request["data"]
	data := storage.CreateNewDataWithTime(dataArgument)
	//TODO add user id to data
	//data["user_id"] = request.Session.Token.UserID
	result := make([]interface{}, 0)

	if err := mongo.Insert(fmt.Sprint(collection), data, &result); err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		ctx.SetBody(common.ConvertInterfaceToJson(result))
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
}

//func save(route *http.Route) {
//	for {
//		select {
//		case <-route.RequestChannel:
//			//request, _ := http.GetRequest(req)
//			//session, _ := auth.GetSession()
//			//
//			//var response http.Response
//			//
//			//collection, _ := request.GetParam("collection")
//			//dataArgument, _ := request.GetParam("data")
//			//if dataArgument == nil {
//			//	fmt.Println(request)
//			//}
//			//data := storage.CreateNewDataWithTime(dataArgument)
//			//data["user_id"] = session.Token.UserID
//			//
//			//insertedData := make([]interface{}, 0)
//			//if err := mongo.Insert(fmt.Sprint(collection), data, &insertedData); err != nil {
//			//	response.SetError(err)
//			//} else {
//			//	response.SetBody(common.ConvertInterfaceToJson(insertedData))
//			//	response.Code = 200
//			//}
//			//
//			//route.ResponseChannel <- &response
//		}
//	}
//}

func savee(request *http.HttpRequest) {
	//var response http.Response
	collection, _ := request.GetParam("collection")
	dataArgument, _ := request.GetParam("data")

	data := storage.CreateNewDataWithTime(dataArgument)
	//if req
	data["user_id"] = request.Session.Token.UserID

	insertedData := make([]interface{}, 0)
	if err := mongo.Insert(fmt.Sprint(collection), data, &insertedData); err != nil {
		request.Response.SetError(err)
	} else {
		request.Response.SetBody(common.ConvertInterfaceToJson(insertedData))
		request.Response.Code = 200
	}
}
