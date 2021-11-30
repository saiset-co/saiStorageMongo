package api

import (
	"encoding/json"
	"fmt"
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/saiset-co/saiStorageMongo/src/sai/auth"
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"github.com/saiset-co/saiStorageMongo/src/sai/db/mongo"
	"github.com/saiset-co/saiStorageMongo/src/sai/network/http"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/routing"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/settings"
)

func Registration() {
	route := &http.Route{
		Name:       "Register",
		Method:     "POST",
		Pattern:    "/register",
		Handler:    registerUser,
		RouteGroup: routing.WithValidation(),
		Validations: []http.Validation{
			{
				Key: "email",
				Rules: []string{
					"!empty",
				},
				Type:     "string",
				Required: true,
			},
			{
				Key: "password",
				Rules: []string{
					"!empty",
				},
				Type:     "string",
				Required: true,
			},
		},
	}
	route.Register()
}

func registerUser(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	user := &auth.User{}
	token := auth.CreateUserToken(user)

	if data, dataExist := request["data"]; dataExist {
		user.Data = data
	}

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

	user.SetCredentials(fmt.Sprint(email), fmt.Sprint(password))

	userData := map[string]interface{}{}
	userData["email"] = email

	if err := mongo.FindOne(auth.UserCollection, userData, nil); err == nil {
		//	ctx.SetStatusCode(err.Code)
		//	ctx.SetBodyString(err.Error())
		//	return
		//} else {
		err = auth.UserAlreadyRegisteredError()
		http.SetErrorResponse(ctx, err)
		return
	}

	roleData := map[string]interface{}{}

	if roleName, roleExists := request["role"]; roleExists {
		roleData["name"] = roleName
	} else {
		roleData["name"] = "user"
	}

	var role auth.Role
	var foundRole map[string]interface{}
	if err := mongo.FindOne(auth.RoleCollection, roleData, &foundRole); err != nil {
		http.SetErrorResponse(ctx, err)
		return
	} else {
		json.Unmarshal(common.ConvertInterfaceToJson(foundRole), &role)
		user.RoleID = role.ID
	}

	insertedUser := []interface{}{}
	if err := mongo.Insert(auth.UserCollection, user, &insertedUser); err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		mongo.Insert(auth.TokenCollection, token, nil)
		http.SetResponse(ctx, common.ConvertInterfaceToJson(insertedUser))
	}
}

func CreateDefaultSuperAdmin(superAdmin settings.SuperAdminConfig) {
	user := &auth.User{}

	user.SetCredentials(superAdmin.Email, superAdmin.Password)

	userData := map[string]interface{}{}
	userData["email"] = user.Email

	if err := mongo.FindOne(auth.UserCollection, userData, nil); err == nil {
		return
	}

	roleData := map[string]interface{}{}
	roleData["name"] = "super_admin"

	var role auth.Role
	var foundRole map[string]interface{}
	if err := mongo.FindOne(auth.RoleCollection, roleData, &foundRole); err != nil {
		return
	}

	json.Unmarshal(common.ConvertInterfaceToJson(foundRole), &role)

	user.RoleID = role.ID
	mongo.Insert(auth.UserCollection, user, nil)
}
