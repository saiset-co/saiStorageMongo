package api

import (
	"saiStorageMongo/src/sai/network/http"
	//"sai/auth"
	"saiStorageMongo/src/sai_storage/routing"
	//"sai/db/mongo"
	"saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"saiStorageMongo/src/sai/db/mongo"
	"saiStorageMongo/src/sai/auth"
	"encoding/json"
	"saiStorageMongo/src/sai/common"
)

func Logout() {
	route := &http.Route{
		Name:       "Logout",
		Method:     "POST",
		Pattern:    "/logout",
		Handler:    logout,
		RouteGroup: routing.WithAuthAndValidation(),
	}
	route.Register()
}

func logout(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	tokenArgument, _ := request["token"]
	tokenData := map[string]interface{}{}
	tokenData["token"] = tokenArgument
	var token auth.Token
	var foundToken interface{}

	err := mongo.FindOne(auth.TokenCollection, tokenData, &foundToken)
	if err != nil {
		err := auth.UserUnauthorizedError()
		http.SetErrorResponse(ctx, err)
		return
	}
	json.Unmarshal(common.ConvertInterfaceToJson(foundToken), &token);
	userData := make(map[string]interface{}, 0)
	userData["user_id"] = token.UserID
	mongo.Remove(auth.TokenCollection, userData, nil)

	ctx.SetBody(common.ConvertInterfaceToJson(token))
	ctx.SetStatusCode(200)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
}
