package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var (
	boldTextStyle   = lipgloss.NewStyle().Bold(true)
	titleStyle      = lipgloss.NewStyle().MarginBottom(1)
	searchStyle     = lipgloss.NewStyle()
	sepStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#c6c6c6"))
	cleanedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7cf4bb"))
	footerStyle     = lipgloss.NewStyle().Margin(1, 0)
	selectedStyle   = boldTextStyle.Copy().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#000000"))
	unselectedStyle = boldTextStyle.Copy()
	deletedStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#ea45b1")).MarginRight(1)
	existingStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#b3f5a3")).MarginRight(1)
)

func viewportContents(m model) string {
	s := ""
	var pathStyle lipgloss.Style
	var sizeStyle lipgloss.Style
	for i, file := range m.directories {
		var boxStyle lipgloss.Style
		if file.Deleted {
			boxStyle = deletedStyle
		} else {
			boxStyle = existingStyle
		}
		if i == m.cursor {
			pathStyle = selectedStyle
			sizeStyle = selectedStyle
		} else {
			pathStyle = unselectedStyle
			sizeStyle = unselectedStyle.Copy().
				Foreground(ColorFromSize(file.Size))
		}

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

	cleaned := cleanedStyle.Render("Cleaned: " + HumanizeBytes(m.bytesSaved))

	s := lipgloss.JoinHorizontal(lipgloss.Top, search, cleaned)

	return titleStyle.Render(s)
}

func (m model) footerView() string {
	return footerStyle.Render(m.help.View(m.keys))
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
