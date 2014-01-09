package blocks

import (
	"fmt"
	"net/http"
)


type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
}


type Controller interface {
}

type ApplicationController struct {
	request *Request
	ViewTemplate string
}


func (ctrl ApplicationController) Dispatch(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
