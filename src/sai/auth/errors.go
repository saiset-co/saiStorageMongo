package auth

import (
	"fmt"
	"github.com/webmakom-com/mycointainer/src/Storage/src/sai/common"
	"net/http"
)

func TokenExpiredError() *common.Error {
	return &common.Error{
		Code: http.StatusUnauthorized,
		What: fmt.Sprintf("Ttoken expired"),
	}
}

func ForbiddenError() *common.Error {
	return &common.Error{
		Code: http.StatusForbidden,
		What: fmt.Sprintf("Permission denied"),
	}
}

func UserAlreadyRegisteredError() *common.Error {
	return &common.Error{
		Code: http.StatusUnauthorized,
		What: fmt.Sprintf("User already registered"),
	}
}

func UserNotRegisteredError() *common.Error {
	return &common.Error{
		Code: http.StatusUnauthorized,
		What: fmt.Sprintf("User not registered"),
	}
}

func UserUnauthorizedError() *common.Error {
	return &common.Error{
		Code: http.StatusUnauthorized,
		What: fmt.Sprintf("User unathorized"),
	}
}
