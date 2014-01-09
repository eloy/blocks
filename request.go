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
	controller Controller      // Controller

	// Response
	contentSet bool             // Has the content already set
	body string                 // Response body
	code int                    // HTML status code

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

func (r *Request) call() {
	// Call the controller for the route
	r.callRequestMethod()

	// Render the view if content isn't already set
	if r.contentSet != true {
		v := NewView(r)
		v.render()
	}

	r.flush()
}

// Instantiate a new controller and call the method
func (r *Request) callRequestMethod() {
	t := r.route.typ
	v := reflect.New(t)
	initializeStruct(t, v.Elem())
	c := v.Interface()//.(*Controller)

	r.controller = c


	// Call the method
	reflect.ValueOf(c).MethodByName(r.route.method).Call(nil)
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
