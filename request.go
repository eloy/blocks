package blocks

import (
	"net/http"
	"reflect"
	"fmt"
)

type Request struct {
	writer http.ResponseWriter   // writer, assigned for the web server
	serverRequest *http.Request // server request, assigned for the web server
	route *Route                // route, assigned by the router.
	scope map[string]string

	// Response
	template string             // Template
	contentSet bool             // Has the content already set
	body string                 // Response body
	code int                    // HTML status code
}

func newRequest(w http.ResponseWriter, r *http.Request) *Request {
	request := new(Request)
	request.writer = w
	request.serverRequest = r
	request.scope = make(map[string]string)
	return request
}

func (this *Request) setRoute(route *Route) {
	this.route = route
	this.template = route.action
}

func (r *Request) setResponse(code int, body string) {
	r.contentSet = true
	r.body = body
	r.code = code
}


func (r *Request) flush() {
	r.writer.WriteHeader(r.code)
	fmt.Fprint(r.writer, r.body)
}


// Execute the request
func (r *Request) call() {

	// Call the controller for the route
	controller := r.callRequestMethod()

	// Render the view if content isn't already set
	if r.contentSet != true {
		v := NewView(r)
		r.setResponse(v.render(controller))
	}

	r.flush()
}

func (r *Request) Path() string {
	return r.serverRequest.URL.Path
}

func (r *Request) Method() string {
	return r.serverRequest.Method
}

// Instantiate a new controller and call the method
func (r *Request) callRequestMethod() Controller {
	t := r.route.controllerT
	v := reflect.New(t)
	initializeStruct(t, v.Elem())
	controller := v.Interface().(Controller)

	controller.setRequest(r)

	// Start the session
	s := newSessionManager(r)
	s.read()
	controller.sessionManager(s)



	// Call the method
	reflect.ValueOf(controller).MethodByName(r.route.action).Call(nil)

	s.save()

	return controller
}
// http://stackoverflow.com/questions/7850140/how-do-you-create-a-new-instance-of-a-struct-from-its-type-at-runtime-in-go
func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch(ft.Type.Kind()){
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
        default:
		}
	}
}
