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

type EntryStatus int

const (
	Alive EntryStatus = iota
	DeletionInProgress
	Deleted
)

type DirEntry struct {
	Path   string
	Size   int64
	Status EntryStatus
}

type FsHandler struct {
	Entries             chan DirEntry
	directoriesToFind   []string
	directoriesToIgnore []string
	searchAll           bool
	shouldFakeEntries   bool
}

func NewFsHandler(c *cli.Context) *FsHandler {
	return &FsHandler{
		Entries:             make(chan DirEntry),
		searchAll:           c.Bool("all"),
		directoriesToFind:   c.Args().Slice(),
		directoriesToIgnore: c.StringSlice("exclude"),
		shouldFakeEntries:   c.Bool("fake-entries"),
	}
}

func (h *FsHandler) Traverse() {
	defer close(h.Entries)

	if h.shouldFakeEntries {
		for _, e := range []DirEntry{
			{Path: "js/vault/node_modules", Size: 2 * gb},
			{Path: "rust/sim/target", Size: 500 * mb},
			{Path: "flutter/app/build", Size: 980 * mb},
			{Path: "flutter/app/.dart_tool", Size: 980 * mb},
			{Path: "js/react_app/node_modules", Size: 1200 * mb},
			{Path: "dart/server/build", Size: 10 * kb},
			{Path: "dart/server/.dart_tool", Size: 10 * kb},
			{Path: "rust/server/target", Size: 600 * mb},
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
	if h.shouldFakeEntries {
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
