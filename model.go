package main

import (
	"log"
	"os"
	"strings"

	"github.com/aarol/goclean/fs"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli"
)

type model struct {
	// Config options
	searchPath  string
	searchDirs  []string
	excludeDirs []string
	searchAll   bool

	viewportReady bool
	height        int
	cursor        int

	searchFinished bool
	bytesSaved     int64

	directories []fs.DirEntry
	sub         chan fs.DirEntry // Will receive directories from fs

	keys     keyMap
	viewport viewport.Model
	spinner  spinner.Model
	help     help.Model
}

func initialModel(c *cli.Context) model {

	// Get search path
	var path string
	var err error
	if c.Bool("home") {
		path, err = os.UserHomeDir()
	} else {
		path, err = os.Getwd()
	}
	if err != nil {
		log.Fatal(err)
	}

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		searchPath:  path,
		searchAll:   c.Bool("all"),
		searchDirs:  c.Args(),
		excludeDirs: strings.Split(c.String("exclude"), " "),

		directories: []fs.DirEntry{},
		sub:         make(chan fs.DirEntry),

		keys:    keys,
		help:    help.New(),
		spinner: s,
	}
}
