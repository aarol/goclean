package main

import (
	"log"
	"os"

	"github.com/karrick/godirwalk"
)

func Traverse(ch chan string) {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defer close(ch)

	// Scan the file tree starting from current working directory
	godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			if directoryEntry.IsDir() && directoryEntry.Name() == os.Args[1] {
				ch <- osPathname
				return godirwalk.SkipThis
			}
			if directoryEntry.Name() == "node_modules" {
				return godirwalk.SkipThis
			}

			return nil
		},
		FollowSymbolicLinks: false,
		Unsorted:            true,
	})
}
