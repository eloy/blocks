package blocks

import (
	"io/ioutil"
	"path"
	"github.com/harlock/peace"
	"net/http"
	"fmt"
)

// Returns the content of the given file
func assetContentFromFile(path string) []byte {
	// TODO: Cache
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}


// Include the content of BlocksRoot/assets/js/application.js
func (this *ApplicationController) JavascriptTag() (buffer string) {
	return javascriptTag()
}

func javascriptTag() (buffer string) {
	source := path.Join(AppRootPath, "assets", "js", "application.js")
	return buildJavascriptTag("/assets", source)
}


func buildJavascriptTag(path string, source string) (buffer string) {
	for _, asset := range peace.Parse(source) {
		buffer += "<script src=\"" + path + "/" + asset + "\" type=\"text/javascript\"></script>\n"
	}
	return buffer
}


func assetsRequestHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		file := r.URL.Path[len(path):]
		peace.WriteResponse(w, file)
	}
}

func EnableAssetsPipeline() {
	peace.AddSource("js", path.Join(AppRootPath, "assets", "js"))
	peace.AddSource("css", path.Join(AppRootPath, "assets", "css"))
	peace.AddVendorSource("js", path.Join(AppRootPath, "vendor", "assets", "js"))
	peace.AddVendorSource("css", path.Join(AppRootPath, "vendor", "assets", "css"))
	http.HandleFunc("/assets/", assetsRequestHandler("/assets/"))
}


// Specs
//----------------------------------------------------------------------
// Include the content of BlocksRoot/spec/application.js
func javascriptTagSpecs() (buffer string) {
	source := path.Join(AppRootPath, "spec", "application.js")
	return buildJavascriptTag("/spec/assets", source)
}


func EnableAssetsPipelineSpecs() {
	peace.AddSource("js", path.Join(AppRootPath, "spec"))
	peace.AddVendorSource("js", path.Join(AppRootPath, "vendor", "assets", "test"))
	http.HandleFunc("/spec/assets/", assetsRequestHandler("/spec/assets/"))
	http.HandleFunc("/specs", specsRunnerRequestHandler)
}

func specsRunnerRequestHandler(w http.ResponseWriter, r *http.Request) {
	template := `
<!DOCTYPE HTML>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <title>Jasmine Spec Runner v2.0.0</title>

  <link rel="stylesheet" type="text/css" href="/assets/jasmine.css">
<!-- Application libs -->
  %s
<!-- Specs libs -->
  %s

</head>

<body>
</body>
</html>
`
	fmt.Fprintf(w, template, javascriptTag(), javascriptTagSpecs())

}
