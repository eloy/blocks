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
	templateName string
	templatePath string
}

func NewView(r *Request) view {
	var v view
	v.templateName = r.template + TEMPLATE_SUFIX
	v.templatePath = path.Join(AppRootPath, VIEWS_DIR, r.route.ControllerName(), v.templateName)
	return v
}

func (v view) render(data interface{}) (int, string) {
	layoutName := DEFAULT_LAYOUT + TEMPLATE_SUFIX
	layoutPath := path.Join(AppRootPath, VIEWS_DIR, LAYOUTS_DIR, layoutName)

	content := string(assetContentFromFile(v.templatePath))
	layout := string(assetContentFromFile(layoutPath))

	body := mustache.RenderInLayout(content, layout, data)
	return 200, body
}
