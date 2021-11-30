package http

import (
	"github.com/saiset-co/saiStorageMongo/src/sai/common"
	"net/http"
)

type Response struct {
	Body  []byte
	Code  int
	Error *common.Error
}

func (resp *Response) SetBody(body []byte) {
	resp.Body = body
}

func (resp *Response) SetBodyString(body string) {
	resp.Body = []byte(body)
}

func (resp *Response) SetError(err *common.Error) {
	resp.Error = err
	resp.SetBodyString(err.Error())
	resp.Code = err.Code
}

func SetResponseError(w *http.ResponseWriter, error2 *common.Error) {
	var response Response
	response.SetError(error2)
	(*w).WriteHeader(response.Code)
	(*w).WriteHeader(response.Code)
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(response.Body)
}
