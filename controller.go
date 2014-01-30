package blocks

import (
	"net/http"
	"encoding/json"
	"github.com/harlock/goson"
	"strconv"
)


type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
}

type Filter func(Controller)

type Controller interface {
	RenderJson(interface{})
	RenderGoson(goson.Goson)
	RedirectTo(string)
	Session() SessionManager
	setRequest(*Request)
	Initialize()
}


type ApplicationController struct {
	// Private fields
	request *Request
}


func (this *ApplicationController) Initialize() {
}

func (this *ApplicationController) setRequest(r *Request) {
	this.request = r
}

// Params
//----------------------------------------------------------------------

func (this *ApplicationController) Param(key string) string {
	// Try first with request params
	param, found := this.request.routeParams[key]
	if found {
		return param
	}
	return this.request.serverRequest.Form.Get(key)
}

func (this *ApplicationController) ParamInt(key string) int {
	value := this.Param(key)
	if value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return 0
}


func (this *ApplicationController) Session() SessionManager {
	return this.request.session
}

func (this *ApplicationController) DecodeJsonBody(model interface{}) error {
	decoder := json.NewDecoder(this.request.serverRequest.Body)
	err := decoder.Decode(model)
	if err != nil {
		panic(err)
	}
	return err
}


// View Helpers
//----------------------------------------------------------------------

func (this *ApplicationController) Forbidden() {
	this.request.setResponse(403, "Forbidden")
}

func (this *ApplicationController) BadRequest() {
	this.request.setResponse(400, "Bad Request")
}


func (this *ApplicationController) RenderJson(object interface{}) {
	json, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	this.request.setResponse(200, string(json))
}

func (this *ApplicationController) RenderGoson(g goson.Goson) {
	json, err := g.ToJson()
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
