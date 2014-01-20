package blocks

import (
	"reflect"
	"log"
	"strings"
	"regexp"
)

type Route struct {
	RouteNode
	controllerT reflect.Type
	method string
	action string
	pathRegExp *regexp.Regexp

}

func newRoute(parent Routable, method string, path string, controller interface{}, action string) *Route {
	r := new(Route)
	r.parent = parent
	r.RouteNode.initialize()
	r.controllerT = reflect.TypeOf(controller)
	r.method = method
	r.setPath(path)
	r.action = action

	return r
}


// Replace variables for regular expressions
// and compile the regexp with the path
func (this *Route) setPath(path string) {

	if strings.HasSuffix(path, "/") {
		path = path + "?"
	}

	this.path = path

	// Replace variables with regexp
	replaced := routeVarsReplaceRegExp.ReplaceAllString(this.Path(), routeVarsReplacement)

	this.pathRegExp = regexp.MustCompile(replaced)
}
var routeVarsReplaceRegExp = regexp.MustCompile(`:([\w]+)`)
const routeVarsReplacement = `(?P<$1>[\w]+)`


func (this *Route) Match(p Pather) bool {
	if p.Method() != this.method {
		return false
	}

	res := this.pathRegExp.MatchString(p.Path())
	// log.Println("Matching request ", p.Path(), ":", p.Method(), "with", this.Path(), "Controller", this.ControllerName(), ":", this.ActionName(), "Res:", res)
	return res
}

// Find Routes
func (this *Route) Search(matches *RoutableCollection, request Pather) {
	// Find first in childrens
	this.findChildrens(matches, request)

	// if childrens don't match, try to with the route itself
	if this.Match(request) {
		matches.add(this)
	}
}


// Return the name lowercase and without controller
// example: HomeController => home
func (r *Route) ControllerName() string {
	return extractControllerName(r.controllerT)
}

func (r *Route) ActionName() string {
	return strings.ToLower(r.action)
}

func (this *Route) Inspect() {
	log.Println(this.Path())
	for _, r := range this.routes {
		r.Inspect()
	}
}
