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
	sessionManager(SessionManager)
}

type ApplicationController struct {
	request *Request
	ViewTemplate string
	Session SessionManager
}

func (this *ApplicationController) setRequest(r *Request) {
	this.request = r
}

func (this *ApplicationController) sessionManager(s SessionManager) {
	this.Session = s
}



func (this *ApplicationController) RenderJson(object interface{}) {
	json, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	this.request.setResponse(200, string(json))
}
