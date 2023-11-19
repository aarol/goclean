package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

type DirEntry struct {
	Path               string
	Size               int64
	Deleted            bool
	DeletionInProgress bool
}

type FsHandler struct {
	Entries             chan DirEntry
	directoriesToFind   []string
	directoriesToIgnore []string
	searchAll           bool
	isDryRun            bool
}

func NewFsHandler(c *cli.Context) *FsHandler {
	return &FsHandler{
		Entries:             make(chan DirEntry),
		searchAll:           c.Bool("all"),
		directoriesToFind:   c.Args().Slice(),
		directoriesToIgnore: c.StringSlice("exclude"),
		isDryRun:            c.Bool("dry-run"),
	}
}

func (h *FsHandler) Traverse() {
	defer close(h.Entries)

	if h.isDryRun {
		for _, e := range []DirEntry{
			{Path: "test/asd", Size: 2 * gb},
			{Path: "test/asd2", Size: 500 * mb},
			{Path: "test/asd3", Size: 10 * kb},
		} {
			time.Sleep(500 * time.Millisecond)
			h.Entries <- e
		}
		return
	}

	// Scan the file tree starting from current working directory
	// Don't return err because it stops the walk
	// Instead just log it and skip the directory
	filepath.WalkDir("", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return filepath.SkipDir
		}
		if d.IsDir() {
			if slices.Contains(h.directoriesToFind, d.Name()) {
				size, err := getDirectorySize(path)
				if err != nil {
					log.Println(err)
					return filepath.SkipDir
				}
				h.Entries <- DirEntry{Path: path, Size: size}
				return filepath.SkipDir
			}
			if slices.Contains(h.directoriesToIgnore, d.Name()) ||
				(strings.HasPrefix(d.Name(), ".") && !h.searchAll) {
				log.Println("Skipping ", path)
				return filepath.SkipDir
			}
		}

		return err
	})
}

func (h *FsHandler) Delete(path string) error {
	if h.isDryRun {
		time.Sleep(500 * time.Millisecond)
		return nil
	}
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
