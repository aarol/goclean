package fs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DirEntry struct {
	Path               string
	Size               int64
	Deleted            bool
	DeletionInProgress bool
}

func Traverse(dir string, searchDirs, ignoreDirs []string, searchAll bool, ch chan DirEntry) {
	defer close(ch)

	// Scan the file tree starting from current working directory
	// Don't return err because it stops the walk
	// Instead just log it and skip the directory
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return filepath.SkipDir
		}
		if d.IsDir() {
			if contains(searchDirs, d.Name()) {
				size, err := getDirectorySize(path)
				if err != nil {
					log.Println(err)
					return filepath.SkipDir
				}
				ch <- DirEntry{Path: path, Size: size}
				return filepath.SkipDir
			}
			if contains(ignoreDirs, d.Name()) ||
				(strings.HasPrefix(d.Name(), ".") && !searchAll) {
				log.Println("Skipping ", path)
				return filepath.SkipDir
			}
		}
		return err
	})
}

func Delete(path string) error {
	return os.RemoveAll(path)
}

// Recursively search every file and return the total sum
func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
