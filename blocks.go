package blocks

import(
	"path"
	"path/filepath"
	"os"
)


var R = newRouter()

// Root path accessor
var root_path string
func RootPath() string {
	return root_path
}


var AppRootPath string

func init() {
	// Set the root path
	root_path, _ = filepath.Abs(filepath.Dir(os.Args[0]))


	// Set the default app root to RootPath/app
	AppRootPath = path.Join(root_path, "app")
}
