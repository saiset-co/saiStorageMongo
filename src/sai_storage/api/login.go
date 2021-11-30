package api

import (
	"github.com/saiset-co/saiStorageMongo/src/sai/auth"
	"github.com/saiset-co/saiStorageMongo/src/sai/network/http"
	//"sai/network/auth"
	"github.com/saiset-co/saiStorageMongo/src/sai/db/mongo"
	//"github.com/saiset-co/saiStorageMongo/src/sai/common"
	//"fmt"
	"encoding/json"
	"fmt"
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/routing"
)

func Login() {
	route := &http.Route{
		Name:       "Login",
		Method:     "POST",
		Pattern:    "/login",
		Handler:    login,
		RouteGroup: routing.WithValidation(),
	}
	route.Register()
}

func login(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	email, emailExist := request["email"]
	if !emailExist {
		err := http.BadRequestError()
		http.SetErrorResponse(ctx, err)
		return
	}

	password, passwordExist := request["password"]
	if !passwordExist {
		err := http.BadRequestError()
		http.SetErrorResponse(ctx, err)
		return
	}

	userData := map[string]interface{}{}
	userData["email"] = email
	userData["password"] = auth.HashUserPassword(fmt.Sprint(password))

	var user auth.User
	var token *auth.Token

	var foundUser map[string]interface{}
	if err := mongo.FindOne(auth.UserCollection, userData, &foundUser); err != nil {
		err = auth.UserNotRegisteredError()
		http.SetErrorResponse(ctx, err)
		return
	}
	json.Unmarshal(common.ConvertInterfaceToJson(foundUser), &user)
	token = auth.CreateUserToken(&user)

	var role auth.Role

	roleData := map[string]interface{}{}
	roleData["_id"] = user.RoleID

	var foundRole map[string]interface{}
	if err := mongo.FindOne(auth.RoleCollection, roleData, &foundRole); err != nil {
		err = mongo.ObjectNotExistsError("role")
		http.SetErrorResponse(ctx, err)
		return
	} else {
		json.Unmarshal(common.ConvertInterfaceToJson(foundRole), &role)
		token = auth.CreateUserToken(&user)

		token.Permissions = role.Permissions
	}

	tokenData := map[string]interface{}{}
	tokenData["user_id"] = token.UserID

	if err := mongo.Remove(auth.TokenCollection, tokenData, nil); err != nil {
		http.SetErrorResponse(ctx, err)
		return
	}
	if err := mongo.Insert(auth.TokenCollection, token, nil); err != nil {
		http.SetErrorResponse(ctx, err)
		return
	}
	auth.AddToken(token)

	http.SetResponse(ctx, common.ConvertInterfaceToJson(token))
}

func SyncTokens() {
	foundTokens := []interface{}{}
	if err := mongo.Find(auth.TokenCollection, nil, nil, &foundTokens); err != nil {
	} else {
		var tokens []auth.Token
		jsonResult, _ := json.Marshal(foundTokens)
		json.Unmarshal(jsonResult, &tokens)
		auth.SyncTokens(tokens)
	}
}
