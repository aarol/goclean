package main

import (
	"log"

	"github.com/aarol/goclean/fs"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type finishedMsg struct{}

// Index of deleted directory
type deletedMsg (int)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenToSearch(m),
		waitForDirectory(m.sub),
		spinner.Tick,
	)
}

func listenToSearch(m model) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan fs.DirEntry)

		go fs.Traverse(m.searchDirs, m.excludeDirs, ch)
		for v := range ch {
			m.sub <- v
		}
		return finishedMsg{}
	}
}

func waitForDirectory(sub chan fs.DirEntry) tea.Cmd {
	return func() tea.Msg {
		e := <-sub
		return e
	}
}

func removeDirectory(index int, path string) tea.Cmd {
	return func() tea.Msg {
		err := fs.Delete(path)
		if err != nil {
			log.Fatal(err)
		}
		return deletedMsg(index)
	}
}

func (m *model) checkPrimaryViewportBounds() {
	top := m.viewport.YOffset
	bottom := m.viewport.Height + top - 1

	if m.cursor < top {
		m.viewport.LineUp(1)
	} else if m.cursor > bottom {
		m.viewport.LineDown(1)
	}

	if m.cursor > len(m.directories)-1 {
		m.cursor = len(m.directories) - 1
	} else if m.cursor < top {
		m.cursor = top
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			m.checkPrimaryViewportBounds()
			m.viewport.SetContent(viewData(m))
		case "down", "j":
			m.cursor++
			m.checkPrimaryViewportBounds()
			m.viewport.SetContent(viewData(m))

		case " ", "enter":
			log.Println("Delete", m.directories[m.cursor].Path)
			return m, removeDirectory(m.cursor, m.directories[m.cursor].Path)
		}

	case tea.WindowSizeMsg:
		if !m.viewportReady {
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			yMargin := headerHeight + footerHeight
			m.viewport = viewport.New(msg.Width, msg.Height-yMargin)
			m.viewport.SetContent(viewData(m))
			m.viewportReady = true
		}

	case fs.DirEntry:
		m.directories = append(m.directories, msg)
		m.viewport.SetContent(viewData(m))
		return m, waitForDirectory(m.sub)

	case finishedMsg:
		m.searchFinished = true

	case deletedMsg:
		m.directories[msg].Deleted = true
		m.viewport.SetContent(viewData(m))

	case spinner.TickMsg:
		if !m.searchFinished {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
