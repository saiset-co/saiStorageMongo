package middleware

import (
	"net/http"
)

func DetectUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//session, _ := auth.GetSession()
		//
		//var response saihttp.Response
		//var user auth.User
		//
		//userData := map[string]interface{}{}
		//userData["id"] = session.Token.UserID
		//
		//var foundUser interface{}
		//if err := mongo.FindOne(auth.UserCollection, userData, &foundUser); err != nil {
		//	response.SetError(auth.UserUnauthorizedError())
		//
		//	tokenData := make(map[string]interface{}, 0)
		//	tokenData["token"] = session.Token
		//	mongo.Remove(auth.TokenCollection, tokenData, nil)
		//	auth.RemoveToken(session.Token.Token)
		//}
		//json.Unmarshal(common.ConvertInterfaceToJson(foundUser), &user);
		//session.User = user
		//
		//if response.Error != nil {
		//
		//	w.WriteHeader(response.Code)
		//	w.Write(response.Body)
		//
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
