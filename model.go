package main

import (
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
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

	directories []DirEntry
	sub         chan DirEntry // Will receive directories from fs

	keys     keyMap
	viewport viewport.Model
	spinner  spinner.Model
	help     help.Model
}

func initialModel(c *cli.Context) model {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)

	return model{
		searchPath:  path,
		searchAll:   c.Bool("all"),
		searchDirs:  c.Args().Slice(),
		excludeDirs: strings.Split(c.String("exclude"), ","),

		directories: []DirEntry{},
		sub:         make(chan DirEntry),

		keys:    keys,
		help:    help.New(),
		spinner: s,
	}
}
