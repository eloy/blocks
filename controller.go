package blocks

import (
	"net/http"
	"encoding/json"
)


type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
}


type Controller interface {
	RenderJson(interface{})
	setRequest(*Request)
}

type ApplicationController struct {
	request *Request
	ViewTemplate string
}

func (this *ApplicationController) setRequest(r *Request) {
	this.request = r
}

func (this *ApplicationController) RenderJson(object interface{}) {
	json, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	this.request.setResponse(200, string(json))
}
