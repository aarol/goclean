package main

import (
	"os"

	"github.com/karrick/godirwalk"
)

func Traverse(ch chan string) {
	// dir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	defer close(ch)

	godirwalk.Walk("C:\\Users\\aarol\\Documents\\Code", &godirwalk.Options{
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			if directoryEntry.Name() == "node_modules" {
				return godirwalk.SkipThis
			}
			if directoryEntry.IsDir() && directoryEntry.Name() == os.Args[1] {
				ch <- osPathname
				return godirwalk.SkipThis
			}

			return nil
		},
		FollowSymbolicLinks: false,
		Unsorted:            true,
	})
}

func contains(slice []string, in string) bool {
	for _, v := range slice {
		if v == in {
			return true
		}
	}
	return false
}
