package blocks

import (
	"reflect"
)

type Pather interface {
	Path() string
	Method() string
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
	r := newRoute(this, "GET", path, controller, action)
	this.addRoute(r)
	return r
}

func (this *RouteNode) Post(path string, controller interface{}, action string) (*Route) {
	r := newRoute(this, "POST", path, controller, action)
	this.addRoute(r)
	return r
}

func (this *RouteNode) Put(path string, controller interface{}, action string) (*Route) {
	r := newRoute(this, "PUT", path, controller, action)
	this.addRoute(r)
	return r
}

func (this *RouteNode) Delete(path string, controller interface{}, action string) (*Route) {
	r := newRoute(this, "DELETE", path, controller, action)
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

	r.Get("/:id/edit", controller, "Edit")
	r.Get("/new", controller, "New")
	r.Get("/:id", controller, "Show")
	r.Get("/", controller, "Index")

	r.Post("/", controller, "Create")
	r.Put("/:id", controller, "Update")
	r.Delete("/:id", controller, "Destroy")

	return r
}

func (this *RouteNode) Member() (Routable) {
	r := newRouteNode()
	p := this.path
	if strings.HasPrefix(p, "/") {
		p = p[1:len(p)]
	}
	r.path = "/:" + p + "_id"
	r.parent = this
	this.addRoute(r)
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

func (this *RouteNode) Search(matches *RoutableCollection, path Pather)  {
	this.findChildrens(matches, path)
}

func (this *RouteNode) findChildrens(matches *RoutableCollection, p Pather)  {
	for _, r := range this.routes {
		r.Search(matches, p)
	}
}

func (this *RouteNode) Inspect() {
	for _, r := range this.routes {
		r.Inspect()
	}
}
