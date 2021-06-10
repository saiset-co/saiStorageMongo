package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/auth"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/common"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/db/mongo"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/network/http"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai_storage/settings"
)

func CheckPermissions(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if settings.Settings.Auth.Enabled {
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

			json.Unmarshal(common.ConvertInterfaceToJson(foundToken), &token)

			url := ctx.Request.URI().Path()
			route, _ := http.API[string(url)]
			if permission, exist := token.Permissions[route.Pattern]; exist {
				for _, reqParam := range permission.Params {
					if param, exist := request[reqParam.Name]; exist {
						if !Can(fmt.Sprint(param), reqParam.Rules, token.UserID) {
							err := auth.ForbiddenError()
							http.SetErrorResponse(ctx, err)
							return
						}
					} else {
						err := auth.ForbiddenError()
						http.SetErrorResponse(ctx, err)
						return
					}
				}
			}
			//
			//if !auth.CurrentSession.User.Can(auth.CurrentSession.RoutePermissions) {
			//	response.SetError(auth.ForbiddenError())
		}
		h(ctx)
	}
}

func Can(param string, rules []string, userId string) bool {
	for _, rule := range rules {
		if rule[:1] == "!" {
			if rule[1:] == param {
				return false
			}
		} else if rule[:1] == "$" {
			switch rule[1:] {
			case "user_id":
				if userId != param {
					return false
				}
				return true
			}
		} else {
			if rule != param {
				return false
			}
		}
	}

	return true
}
