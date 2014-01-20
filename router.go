package blocks

import (
	"net/http"
	"reflect"
	"log"
	"strings"
	"path"
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

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	log.Println("Serving Request", r.Method, r.URL, r.PostForm)
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
	} else {
		// If not found, use a fileServer
		public := path.Join(AppRootPath, "public")
		server := http.FileServer(http.Dir(public))
		server.ServeHTTP(w,r)
	}

}

// Clean the router. Used in tests
func (this *Router) Reset() {
	this.rootRoute = nil
	this.RouteNode.initialize()
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
	matches := new(RoutableCollection)
	this.findChildrens(matches, path)

	switch len(matches.routes) {
	case 0:
		return nil, false
	case 1:
		return matches.routes[0], true
	default:
		size := 0
		var c *Route
		for _, route := range(matches.routes) {
			if current_size := len(route.Path()); current_size > size {
				size = current_size
				c = route
			}
		}
		if size != 0 {
			return c, true
		}
	}
	return nil, false
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
