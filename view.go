package blocks

import (
	"github.com/hoisie/mustache"
	"path"
)

const TEMPLATE_SUFIX = ".mustache"
const VIEWS_DIR = "views"
const LAYOUTS_DIR = "layouts"
const DEFAULT_LAYOUT = "application"

type view struct {
	request *Request
	templateName string
	templatePath string
}

func NewView(r *Request) view {
	var v view
	v.request = r
	v.templateName = r.route.ActionName() + TEMPLATE_SUFIX
	v.templatePath = path.Join(AppRootPath, VIEWS_DIR, r.route.ControllerName(), v.templateName)
	return v
}

func (v view) render() {
	layoutName := DEFAULT_LAYOUT + TEMPLATE_SUFIX
	layoutPath := path.Join(AppRootPath, VIEWS_DIR, LAYOUTS_DIR, layoutName)

	content := string(assetContentFromFile(v.templatePath))
	layout := string(assetContentFromFile(layoutPath))

	v.request.body = mustache.RenderInLayout(content, layout, v.request.controller)
	v.request.code = 200
}
