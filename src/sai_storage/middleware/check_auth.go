package middleware

import (
	"encoding/json"
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/kirillbeldyaga/fasthttp"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/auth"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/common"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/db/mongo"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/network/http"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai_storage/settings"
)

func CheckAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if settings.Settings.Auth.Enabled {
			var request map[string]interface{}
			json.Unmarshal(ctx.PostBody(), &request)

			tokenArgument, exist := request["token"]
			if !exist {
				err := http.BadRequestError()
				http.SetErrorResponse(ctx, err)
				return
			}
			//userId := request["user_id"]

			//	request, _ := http2.GetRequest(r)
			//	session, _ := auth.GetSession()
			//
			//	tokenArgument, exist := request.GetParam("token")
			//
			//	if !exist {
			//		return
			//	}
			//
			tokenData := map[string]interface{}{}
			tokenData["token"] = tokenArgument
			//
			var token auth.Token
			//	if cachedToken, exist := auth.Tokens[fmt.Sprint(tokenArgument)]; !exist {
			var foundToken interface{}

			err := mongo.FindOne(auth.TokenCollection, tokenData, &foundToken)
			if err != nil {
				err := auth.UserUnauthorizedError()
				http.SetErrorResponse(ctx, err)
				return
			}
			json.Unmarshal(common.ConvertInterfaceToJson(foundToken), &token)

			//		auth.AddToken(&token)
			//	} else {
			//		token = cachedToken
			//	}
			if expired := token.VerifyToken(); expired {
				mongo.Remove(auth.TokenCollection, tokenData, nil)
				err := auth.TokenExpiredError()
				http.SetErrorResponse(ctx, err)

				return
			}
			//
			//	session.Token = token
		}
		h(ctx)
	}
}
