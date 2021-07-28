package common

import (
	"encoding/json"
)

func ConvertInterfaceToJson(obj interface{}) []byte {
	jsonResult, _ := json.Marshal(obj)
	return jsonResult
}

func ConvertJsonToInterface(encodedJson string) (interface{}, *Error) {
	var inter interface{}

	err := json.Unmarshal([]byte(encodedJson), &inter)
	if err != nil {
		return nil, InvalidArgumentError(encodedJson)
	}

	return inter, nil
}
