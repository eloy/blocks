package blocks

import (
	"net/http"
	"reflect"
	"log"
	"strings"
	"regexp"
	"runtime/debug"
)

type Router struct {
	RouteNode
	rootRoute *Route
}

func newRouter() *Router {
	r := new(Router)
	r.RouteNode.initialize()
	return r
}

func (this *Router) Path() string {
	return "/"
}

// Server the requests
func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := newRequest(w, r)

	defer func() {
		if err := recover(); err != nil {
			log.Println("ERROR:", err)
			http.Error(w, (err.(error)).Error(), http.StatusInternalServerError)
			// fmt.Fprintf(w, "ERROR %v", err)
			debug.PrintStack()
		}
	}()

	route, found := this.Find(request)
	if found {
		request.setRoute(route)
		request.call()
	}

	// TODO Not Found

}


func (this *Router) Root(controller interface{}, action string) (*Route) {
	r := newRoute(this, "GET", "/", controller, action)
	this.rootRoute = r
	return r
}

func (this *Router) Find(path Pather) (*Route, bool)  {
	if path.Path() == "/" {
		return this.rootRoute, true
	}
	return this.findChildrens(path)
}


var cleanupPathRegexp = regexp.MustCompile("/+")
func httpJoin(strs ...string) string {
	path := strings.Join(strs, "/")
	return cleanupPathRegexp.ReplaceAllString(path, "/")
}


func extractControllerName(t reflect.Type) string {
	name := t.Name()
	return strings.ToLower(name[:len(name) - 10])

}
