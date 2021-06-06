package http

import (
	"fmt"
	"saiStorageMongo/src/sai/common"
	"net/http"
)

func RouteNotFindError(path string) *common.Error {
	return &common.Error{
		Code: http.StatusNotFound,
		What: fmt.Sprintf("Route %s not found", path),
	}
}

func MethodNotAllowedError(method string) *common.Error {
	return &common.Error{
		Code: http.StatusMethodNotAllowed,
		What: fmt.Sprintf("Method %s not allowed", method),
	}
}

func BadRequestError() *common.Error {
	return &common.Error{
		Code: http.StatusBadRequest,
		What: fmt.Sprintf("Bad request"),
	}
}
