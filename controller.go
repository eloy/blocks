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
	Session SessionManager

	// Private fields
	request *Request
}

func (this *ApplicationController) setRequest(r *Request) {
	this.request = r
}

func (this *ApplicationController) Params(key string) string {
	return this.request.serverRequest.FormValue(key)
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

func (this *ApplicationController) RedirectTo(url string) {
	this.request.writer.Header()["location"] = []string{url}
	this.request.setResponse(302, "")
}


func (this *ApplicationController) RenderTemplate(template string) {
	this.request.template = template
}
