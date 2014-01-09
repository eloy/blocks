package blocks

import (
	"html/template"
	"path"
	"strings"
)

const TEMPLATE_SUFIX = ".html.go"

type view struct {
}


func (v view) render(r *Request) {
	method := strings.ToLower(r.route.method)
	template_name := method + TEMPLATE_SUFIX
	controller_name := strings.ToLower(r.route.ControllerName())

	template_path := path.Join(AppRootPath, "views", controller_name, template_name)

	t, err := template.New(template_name).ParseFiles(template_path)
	if err != nil {
		panic(err)
	}

	err = t.Execute(r.writer, r.controller)
	if err != nil {
		panic(err)
	}

}
