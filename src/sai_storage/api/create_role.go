package api

import (
	"github.com/saiset-co/saiStorageMongo/src/sai/auth"
	"github.com/saiset-co/saiStorageMongo/src/sai/db/mongo"
	"github.com/saiset-co/saiStorageMongo/src/sai/network/http"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/routing"
	//"encoding/json"
	"github.com/saiset-co/saiStorageMongo/src/sai_storage/settings"
	//"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"encoding/json"
	"github.com/saiset-co/saiStorageMongo/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
)

func CreateRole() {
	route := &http.Route{
		Name:       "CreateRole",
		Method:     "POST",
		Pattern:    "/role_create",
		Handler:    createRole,
		RouteGroup: routing.WithAuthAndValidation(),
		Validations: []http.Validation{
			{
				Key:      "data",
				Type:     "object",
				Required: true,
			},
		},
	}
	route.Register()
}

func createRole(ctx *fasthttp.RequestCtx) {
	var request map[string]interface{}
	json.Unmarshal(ctx.PostBody(), &request)

	data := request["data"]

	role := &auth.Role{}
	if roleData, err := json.Marshal(data); err != nil {
	} else {
		json.Unmarshal([]byte(roleData), &role)
	}

	roleData := map[string]interface{}{}
	roleData["name"] = role.Name

	if err := mongo.FindOne(auth.RoleCollection, roleData, nil); err == nil {
		err = mongo.ObjectAlreadyExistsError("role")
		http.SetErrorResponse(ctx, err)
		return
	}

	createdRole := []interface{}{}
	if err := mongo.Insert(auth.RoleCollection, role, &createdRole); err != nil {
		http.SetErrorResponse(ctx, err)
	} else {
		ctx.SetBody(common.ConvertInterfaceToJson(createdRole))
		ctx.SetStatusCode(200)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	}
}

func CreateDefaultRoles(roles []settings.RoleConfig) {
	for _, roleConfig := range roles {
		roleData := map[string]interface{}{}
		roleData["name"] = roleConfig.Name

		if err := mongo.FindOne(auth.RoleCollection, roleData, nil); err == nil {
			return
		}

		role := auth.CreateNewRole()
		role.Name = roleConfig.Name
		role.Permissions = map[string]auth.Permission{}

		for _, permission := range roleConfig.Permissions {
			params := []auth.Param{}

			for _, param := range permission.Params {
				params = append(params, auth.Param{
					Name:  param.Name,
					Rules: param.Rules,
				})
			}
			role.Permissions[permission.URL] = auth.Permission{
				URL:    permission.URL,
				Rights: permission.Rights,
				Params: params,
			}
		}

		mongo.Insert(auth.RoleCollection, role, nil)
	}
}
