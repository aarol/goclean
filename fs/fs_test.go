package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPopulateDirectories(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(filepath.Dir(wd), "example", "output.go")
	os.MkdirAll(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	path = filepath.Join(filepath.Dir(wd), "nested", "example", "output.go")
	os.MkdirAll(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
