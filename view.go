package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var (
	headerStyle        = lipgloss.NewStyle().MarginBottom(1)
	searchStyle        = lipgloss.NewStyle()
	cleanedStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("121"))
	sepStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("249"))

	boldTextStyle = lipgloss.NewStyle().Bold(true)

	selectedPathStyle   = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#000"))
	unselectedPathStyle = lipgloss.NewStyle()

	iconStyle         = lipgloss.NewStyle().MarginRight(1)
	deletedIconStyle  = iconStyle.Copy().Background(lipgloss.Color("205"))
	deletingIconStyle = iconStyle.Copy().Background(lipgloss.Color("214"))
	aliveIconStyle    = iconStyle.Copy().Background(lipgloss.Color("121"))
)

// Updates viewport with the data in model
func (m *model) updateViewport() {
	var s strings.Builder
	var pathStyle lipgloss.Style
	var sizeStyle lipgloss.Style
	var boxStyle lipgloss.Style

	for i, file := range m.directories {
		// Box on the left of path
		switch file.Status {
		case Deleted:
			boxStyle = deletedIconStyle
		case DeletionInProgress:
			boxStyle = deletingIconStyle
		case Alive:
			boxStyle = aliveIconStyle
		case Error:
			boxStyle = lipgloss.NewStyle().MarginRight(1).Background(lipgloss.Color("#ff0000"))
		}

		if i == m.cursor {
			pathStyle = selectedPathStyle
			sizeStyle = selectedPathStyle.Bold(true)
		} else {
			pathStyle = unselectedPathStyle
			sizeStyle = unselectedPathStyle.Copy().
				Bold(true).
				Foreground(ColorFromSize(file.Size))
		}

		box := boxStyle.Render(" ")

		fileSize := sizeStyle.Render(HumanizeBytes(file.Size))

		directoryWidth := m.viewport.Width - lipgloss.Width(fileSize) - lipgloss.Width(box)

		var directoryItem string

		if file.Status == Error {
			directoryItem = pathStyle.Copy().
				Width(directoryWidth).
				Render(truncate.StringWithTail(
					file.ErrorMsg, uint(directoryWidth), "..."),
				)
		} else {
			directoryItem = pathStyle.Copy().
				Width(directoryWidth).
				Render(
					truncate.StringWithTail(
						file.Path, uint(directoryWidth), "..."),
				)
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, box, directoryItem, fileSize)

		fmt.Fprintln(&s, row)
	}
	m.viewport.SetContent(s.String())
}

func (m model) headerView() string {
	var search string
	if !m.searchFinished {
		dirs := strings.Join(m.FsHandler.directoriesToFind, ", ")
		search = searchStyle.Render(fmt.Sprintf("%s Scanning directories %s • ", m.spinner.View(), dirs))
	} else {
		s := fmt.Sprintf("Found %d directories • ", len(m.directories))
		search = searchStyle.Copy().Inherit(sepStyle).Render(s)
	}

	cleaned := cleanedStatusStyle.Render("Cleaned: " + HumanizeBytes(m.bytesSaved))

	s := lipgloss.JoinHorizontal(lipgloss.Top, search, cleaned)

	s += boldTextStyle.Render("\nPress DELETE/ENTER to delete directory")

	return headerStyle.Render(s)
}

var footerStyle = lipgloss.NewStyle().Margin(1, 0)

func (m model) footerView() string {
	return footerStyle.Render(m.help.View(m.keys))
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
