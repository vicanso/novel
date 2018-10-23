package asset

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

var adminAsset *Asset
var webAsset *Asset

func init() {
	adminAsset = &Asset{
		box: packr.NewBox("../admin/dist"),
	}
	webAsset = &Asset{
		box: packr.NewBox("../assets"),
	}
}

// GetAdminAsset get admin asset
func GetAdminAsset() *Asset {
	return adminAsset
}

// GetWebAsset get web asset
func GetWebAsset() *Asset {
	return webAsset
}

// Asset asset
type Asset struct {
	box packr.Box
}

// Open open the file
func (a *Asset) Open(filename string) (http.File, error) {
	return a.box.Open(filename)
}

// Get the the data of file
func (a *Asset) Get(filename string) []byte {
	return a.box.Bytes(filename)
}

// Exists check the file exists
func (a *Asset) Exists(filename string) bool {
	return a.box.Has(filename)
}
