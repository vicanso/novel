package asset

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

var box packr.Box

func init() {
	box = packr.NewBox("../assets")
}

// Open open the file
func Open(filename string) (http.File, error) {
	return box.Open(filename)
}

// Get the the data of file
func Get(filename string) []byte {
	return box.Bytes(filename)
}

// Exists check the file exists
func Exists(filename string) bool {
	return box.Has(filename)
}
