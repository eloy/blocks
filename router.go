package blocks

import (
	"net/http"
	"reflect"
	"log"
)


type Pathed interface {
	Path() string
}

type Route struct {
	typ reflect.Type
	parent Pathed
	path string
	method string
}

type Router struct {
}

var routes = make([]*Route, 0)

func addRoute(r *Route) {
	routes = append(routes, r)
}


func findRoute(path string) *Route {
	for _, r := range routes {
		if path == r.path {
			return r
		}
	}
	panic("Route not found")
}


func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := &Request{writer: w, serverRequest: r}
	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR:", err)
			http.Error(w, (err.(error)).Error(), http.StatusInternalServerError)
			// fmt.Fprintf(w, "ERROR %v", err)
		}
	}()

	request.route = findRoute(r.URL.Path)
	request.call()


}


func (this *Router) Get(path string, controller interface{}, method string) (*Route) {
	r := new(Route)
	r.typ = reflect.TypeOf(controller)
	r.path = path
	r.method = method
	addRoute(r)
	return r

}



// func (this *Router) Resources(controller interface{}) (*Route) {
// 	r := new(Route)
// 	r.typ = reflect.TypeOf(controller)
// 	addRoute(r)
// 	return r
// }


func (r *Route) ControllerName() string {
	name := r.typ.Name()
	return r.typ.Name()[:len(name) - 10]
}
