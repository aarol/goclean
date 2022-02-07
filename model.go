package main

import (
	"strings"

	"github.com/aarol/goclean/fs"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli"
)

type model struct {
	viewportReady  bool
	searchFinished bool
	searchDirs     []string
	excludeDirs    []string
	directories    []fs.DirEntry
	sub            chan fs.DirEntry
	viewport       viewport.Model
	spinner        spinner.Model
	cursor         int
}

func initialModel(c *cli.Context) model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		directories: []fs.DirEntry{},
		sub:         make(chan fs.DirEntry),
		searchDirs:  c.Args(),
		excludeDirs: strings.Split(c.String("exclude"), " "),
		spinner:     s,
	}
}
