package fs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type DirEntry struct {
	Path    string
	Size    int64
	Deleted bool
}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}

func Traverse(searchDirs, ignoreDirs []string, ch chan DirEntry) {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// dir := "C:\\Users\\aarol\\Documents\\Code"
	defer close(ch)

	// Scan the file tree starting from current working directory
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if contains(d.Name(), searchDirs) {
				size, err := getDirectorySize(path)
				if err != nil {
					return err
				}
				ch <- DirEntry{Path: path, Size: size}
				return filepath.SkipDir
			}
			if contains(d.Name(), ignoreDirs) {
				return filepath.SkipDir
			}
		}
		return err
	})
}

func Delete(path string) error {
	return os.RemoveAll(path)
}

func contains(e string, s []string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
