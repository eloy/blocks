package blocks

import (
	"reflect"
)

type Pather interface {
	Path() string
}

type Routable interface {
	Path() string
	Find(Pather) (*Route, bool)
	Match(Pather) bool
	Inspect()

	Get(string, interface{}, string) *Route
	Resources(interface{}) Routable
	Namespace(string) Routable
}


type RouteNode struct {
	parent Routable
	path string
	routes []Routable
}

func newRouteNode() *RouteNode {
	r := new(RouteNode)
	r.initialize()
	return r
}

func (this *RouteNode) initialize() {
	this.routes = make([]Routable, 0)
}

func (this *RouteNode) Get(path string, controller interface{}, action string) (*Route) {
	r := newRoute(this, path, controller, action)
	this.addRoute(r)
	return r
}

func (this *RouteNode) Namespace(path string) (Routable) {
	r := newRouteNode()
	r.path = path
	r.parent = this
	this.addRoute(r)
	return r
}

func (this *RouteNode) Resources(controller interface{}) (Routable) {
	controllerT :=  reflect.TypeOf(controller)
	path := "/" + extractControllerName(controllerT)

	r := newRouteNode()
	r.path = path
	r.parent = this
	this.addRoute(r)


	// Index
	r.Get("/", controller, "Index")
	// Edit
	r.Get("/edit/:Id", controller, "Edit")

	return r
}

func (this *RouteNode) addRoute(r Routable) {
	this.routes = append(this.routes, r)
}

func (this *RouteNode) Path() string {
	parentPath := ""
	if this.parent != nil {
		parentPath = this.parent.Path()
	}

	return httpJoin(parentPath, this.path)
}

// Always return false
// See Route.Match for implementation
func (this *RouteNode) Match(p Pather) bool {
	return false
}

func (this *RouteNode) Find(path Pather) (*Route, bool)  {
	return this.findChildrens(path)
}

func (this *RouteNode) findChildrens(p Pather) (*Route, bool)  {
	for _, r := range this.routes {
		route, found := r.Find(p)
		if found {
			return route, true
		}
	}
	return nil, false
}

func (this *RouteNode) Inspect() {
	for _, r := range this.routes {
		r.Inspect()
	}
}
