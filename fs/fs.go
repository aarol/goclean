package fs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
)

type DirEntry struct {
	Path    string
	Size    int64
	Deleted bool
}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err

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
	godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, entry *godirwalk.Dirent) error {
			if entry.IsDir() {
				if contains(entry.Name(), searchDirs) {
					size, err := getDirectorySize(osPathname)
					if err != nil {
						return err
					}
					ch <- DirEntry{Path: osPathname, Size: size}
					return godirwalk.SkipThis
				}
				if contains(entry.Name(), ignoreDirs) {
					return godirwalk.SkipThis
				}
			}

			return nil
		},
		FollowSymbolicLinks: false,
		Unsorted:            true,
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
