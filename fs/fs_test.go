package fs

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func TestPopulateDirectories(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(filepath.Dir(wd), "example")
	err = os.MkdirAll(path, os.ModeAppend)
	if err != nil {
		t.Fatal(err)
	}
	err = copyFileContents("../CascadiaCode.zip", filepath.Join(path, "CascadiaCode.zip"))
	if err != nil {
		t.Fatal(err)
	}

	path = filepath.Join(filepath.Dir(wd), "nested", "example")
	os.MkdirAll(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	path = filepath.Join(filepath.Dir(wd), "other", "example")
	os.MkdirAll(path, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
