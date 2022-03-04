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

	selectedPathStyle   = boldTextStyle.Copy().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#000"))
	unselectedPathStyle = boldTextStyle.Copy()

	deletedIconStyle  = lipgloss.NewStyle().Background(lipgloss.Color("205")).MarginRight(1)
	existingIconStyle = lipgloss.NewStyle().Background(lipgloss.Color("121")).MarginRight(1)

	footerStyle = lipgloss.NewStyle().Margin(1, 0)
)

func viewportContents(m model) string {
	s := ""
	var pathStyle lipgloss.Style
	var sizeStyle lipgloss.Style
	for i, file := range m.directories {
		var boxStyle lipgloss.Style
		// Box on the left of path
		if file.Deleted {
			boxStyle = deletedIconStyle
		} else {
			boxStyle = existingIconStyle
		}
		if i == m.cursor {
			pathStyle = selectedPathStyle
			sizeStyle = selectedPathStyle
		} else {
			pathStyle = unselectedPathStyle
			sizeStyle = unselectedPathStyle.Copy().
				Foreground(ColorFromSize(file.Size))
		}
		// Set boldness to true if file not deleted
		pathStyle.Bold(!file.Deleted)

		box := boxStyle.Render(" ")

		fileSize := sizeStyle.Render(HumanizeBytes(file.Size))

		directoryWidth := m.viewport.Width - lipgloss.Width(fileSize) - lipgloss.Width(box)

		directoryItem :=
			pathStyle.Copy().
				Width(directoryWidth).
				Render(
					truncate.StringWithTail(
						file.Path, uint(directoryWidth), "..."),
				)

		row := lipgloss.JoinHorizontal(lipgloss.Top, box, directoryItem, fileSize)

		s += fmt.Sprintln(row)
	}
	return s
}

func (m model) headerView() string {
	var search string
	if !m.searchFinished {
		dirs := strings.Join(m.searchDirs, ", ")
		search = searchStyle.Render(fmt.Sprintf("%s Scanning directories %s • ", m.spinner.View(), dirs))
	} else {
		s := fmt.Sprintf("Found %d directories • ", len(m.directories))
		search = searchStyle.Copy().Inherit(sepStyle).Render(s)
	}

	cleaned := cleanedStatusStyle.Render("Cleaned: " + HumanizeBytes(m.bytesSaved))

	s := lipgloss.JoinHorizontal(lipgloss.Top, search, cleaned)

	s += boldTextStyle.Render("\nPress SPACE to delete directory")

	return headerStyle.Render(s)
}

func (m model) footerView() string {
	return footerStyle.Render(m.help.View(m.keys))
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
