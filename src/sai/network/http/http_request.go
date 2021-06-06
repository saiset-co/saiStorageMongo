package http

import (
	"net/url"
	"net/http"
	"saiStorageMongo/src/sai/common"
	"saiStorageMongo/src/sai/auth"
	"encoding/json"
)

//type Param interface{}

type HttpRequest struct {
	Query           url.Values
	Body            map[string]interface{}
	Params          map[string]interface{}
	Route           *Route
	Session         *auth.Session
	Response        *Response
	ResponseChannel chan *Response
}

var (
//	CurrentRequest = &HttpRequest{
//		Body:   make(map[string]interface{}),
//		Params: make(map[string]interface{}),
//	}
//RequestMutex = &sync.Mutex{}
)

func CreateRequest() HttpRequest {
	return HttpRequest{
		Body:            make(map[string]interface{}),
		Params:          make(map[string]interface{}),
		Session:         auth.CreateSession(),
		ResponseChannel: make(chan *Response),
		Response:        &Response{},
	}
}

func GetRequest(req *http.Request, route *Route) (*HttpRequest, *common.Error) {
	request := CreateRequest()
	request.Query = req.URL.Query()
	request.Route = route

	if req.Method == "POST" {
		err := json.NewDecoder(req.Body).Decode(&request.Body)
		if err != nil {
			return nil, BadRequestError()
		}
	}

	request.ParseParams()

	return &request, nil
}

func (httpRequest *HttpRequest) ParseParams() {
	if len(httpRequest.Body) > 0 {
		httpRequest.Params = httpRequest.Body
	}
	for k, v := range httpRequest.Query {
		httpRequest.Params[k] = v[0]
	}
}

func (httpRequest *HttpRequest) GetParam(key string) (interface{}, bool) {
	value, exist := httpRequest.Params[key]

	return value, exist
}
