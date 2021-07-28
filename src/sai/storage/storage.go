package storage

import (
	"github.com/webmakom-com/mycointainer/src/Storage/src/github.com/segmentio/ksuid"
	"time"
)

type Data map[string]interface{}

func CreateNewDataWithTime(newData interface{}) map[string]interface{} {
	//var data map[string]interface{}
	//if newData == nil {
	//	data = map[string]interface{}{}
	//} else {
		data := newData.(map[string]interface{})
	//}
	// data["cr_time"] = time.Now().Format("2006-01-02 15:04:05")
	data["cr_time"] = time.Now().UnixNano()
	data["internal_id"] = CreateDocumentID()
	return data
}

func CreateNewData(newData interface{}) map[string]interface{} {
	return newData.(map[string]interface{})
}

func CopyToData(data map[string]interface{}, newData interface{}) map[string]interface{} {
	for k, v := range newData.(map[string]interface{}) {
		data[k] = v
	}
	return data
}

func CreateDocumentID() string {
	return ksuid.New().String()
}
