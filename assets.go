package blocks

import (
	"io/ioutil"
)

func assetContentFromFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}
