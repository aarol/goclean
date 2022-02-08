package main

import (
	"log"

	"github.com/aarol/goclean/fs"
	"github.com/charmbracelet/bubbles/key"
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
		m.spinner.Tick,
	)
}

// Will return finishedMsg when completed search
func listenToSearch(m model) tea.Cmd {
	return func() tea.Msg {
		// Channel receives results, then immediately sends them to model channel
		ch := make(chan fs.DirEntry)

		go fs.Traverse(m.searchDirs, m.excludeDirs, ch)
		for v := range ch {
			m.sub <- v
		}
		return finishedMsg{}
	}
}

// Receives directories from model channel
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

// Keeps cursor visible on viewport, and inside bounds
func (m *model) checkPrimaryViewportBounds() {
	// Vertical scroll position of viewport
	top := m.viewport.YOffset
	// Last position visible in viewport
	bottom := m.viewport.Height + top - 1

	if m.cursor < top {
		m.viewport.LineUp(1)
	} else if m.cursor > bottom {
		// For some reason it still scrolls down once at the last index
		// Needs further investigating
		m.viewport.LineDown(1)
	}

	// Keep cursor within bounds
	if m.cursor > len(m.directories)-1 {
		m.cursor = len(m.directories) - 1
	} else if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			m.cursor--
			m.checkPrimaryViewportBounds()
			m.viewport.SetContent(viewportContents(m))

		case key.Matches(msg, m.keys.Down):
			m.cursor++
			m.checkPrimaryViewportBounds()
			m.viewport.SetContent(viewportContents(m))

		case key.Matches(msg, m.keys.Delete):
			log.Println("Delete", m.directories[m.cursor].Path)
			return m, removeDirectory(m.cursor, m.directories[m.cursor].Path)

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

			// Recalculate viewport height, since help is now occupying more space
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			yMargin := headerHeight + footerHeight
			m.viewport.Height = m.height - yMargin
		}

	case tea.WindowSizeMsg:
		// Terminal dimensions are sent asynchronously
		if !m.viewportReady {
			m.height = msg.Height
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			yMargin := headerHeight + footerHeight
			// Viewport occupies maximum height
			m.viewport = viewport.New(msg.Width, msg.Height-yMargin)
			m.viewport.SetContent(viewportContents(m))
			// Help needs the width for truncating
			// Otherwise renders nothing
			m.help.Width = msg.Width

			m.viewportReady = true
		}

	case fs.DirEntry:
		m.directories = append(m.directories, msg)
		m.viewport.SetContent(viewportContents(m))
		return m, waitForDirectory(m.sub)

	case finishedMsg:
		m.searchFinished = true

	case deletedMsg:
		m.directories[msg].Deleted = true
		m.bytesSaved += m.directories[msg].Size
		m.viewport.SetContent(viewportContents(m))

	case spinner.TickMsg:
		if !m.searchFinished {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
