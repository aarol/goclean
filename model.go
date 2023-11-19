package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
)

type model struct {
	FsHandler *FsHandler

	viewportReady bool
	height        int
	cursor        int

	searchFinished bool
	bytesSaved     int64

	directories []DirEntry

	keys     keyMap
	viewport viewport.Model
	spinner  spinner.Model
	help     help.Model
}

func initialModel(c *cli.Context) model {
	s := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
	)

	fsHandler := NewFsHandler(c)

	return model{
		FsHandler:   fsHandler,
		directories: []DirEntry{},

		keys:    keys,
		help:    help.New(),
		spinner: s,
	}
}
