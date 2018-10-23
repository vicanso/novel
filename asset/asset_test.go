package asset

import (
	"testing"
)

func TestAsset(t *testing.T) {
	filename := "index.html"
	t.Run("open", func(t *testing.T) {
		f, err := adminAsset.Open(filename)
		if err != nil {
			t.Fatalf("open fail, %v", err)
		}
		fi, err := f.Stat()
		if err != nil {
			t.Fatalf("stat fail, %v", err)
		}
		if fi.Name() != filename {
			t.Fatalf("get stat info fail")
		}
	})

	t.Run("get", func(t *testing.T) {
		buf := adminAsset.Get(filename)
		if len(buf) == 0 {
			t.Fatalf("get file data fail")
		}
	})

	t.Run("exists", func(t *testing.T) {
		if !adminAsset.Exists(filename) {
			t.Fatalf("check exists fail")
		}
	})
}
