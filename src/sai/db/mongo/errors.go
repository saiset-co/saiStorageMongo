package mongo

import (
	"saiStorageMongo/src/sai/common"
	"net/http"
	"fmt"
)

func ObjectAlreadyExistsError(objectName string) *common.Error {
	return &common.Error{
		Code: http.StatusBadRequest,
		What: fmt.Sprintf("Object %s already exist", objectName),
	}
}

func ObjectNotExistsError(objectName string) *common.Error {
	return &common.Error{
		Code: http.StatusBadRequest,
		What: fmt.Sprintf("Object %s not exist", objectName),
	}
}

func InvalidObjectIdError(objectId string) *common.Error {
	return &common.Error{
		Code: http.StatusBadRequest,
		What: fmt.Sprintf("Invalid object id %s", objectId),
	}
}

func MongoDBError(b error) *common.Error {
	return &common.Error{
		Code: http.StatusInternalServerError,
		What: fmt.Sprintf(b.Error()),
	}
}
