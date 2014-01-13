package blocks

import (
	"io/ioutil"
	"path"
	"github.com/harlock/peace"
	"net/http"
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
	source := path.Join(AppRootPath, "assets", "js", "application.js")
	for _, asset := range peace.Parse(source) {
		buffer += "<script src=\"/assets/" + asset + "\" type=\"text/javascript\"></script>\n"
	}
	return buffer
}


func assetsRequestHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/assets/"):]
	peace.WriteResponse(w, file)
}

func EnableAssetsPipeline() {
	peace.AddSource("js", path.Join(AppRootPath, "assets", "js"))
	peace.AddVendorSource("js", path.Join(AppRootPath, "vendor", "assets", "js"))
	http.HandleFunc("/assets/", assetsRequestHandler)
}
