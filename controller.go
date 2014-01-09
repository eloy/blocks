package blocks

import (
	"fmt"
	"net/http"
)


type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
}


type Controller struct {
	request *Request
}


func (ctrl Controller) Dispatch(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
