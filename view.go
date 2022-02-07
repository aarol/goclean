package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var (
	titleStyle      = lipgloss.NewStyle().MarginBottom(1)
	footerStyle     = lipgloss.NewStyle().Margin(1, 0)
	boldTextStyle   = lipgloss.NewStyle().Bold(true)
	selectedStyle   = boldTextStyle.Copy().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#000000"))
	unselectedStyle = boldTextStyle.Copy()
	deletedStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#ea45b1")).MarginRight(1)
	existingStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#b3f5a3")).MarginRight(1)
)

func viewData(m model) string {
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
				// Foreground(ColorFromSize(file.Size))
				Foreground(lipgloss.Color("#000000"))
		}

		box := boxStyle.Render(" ")

		log.Println(sizeStyle.Render(HumanizeBytes(file.Size)))
		fileSize := sizeStyle.Render(HumanizeBytes(file.Size))

		directoryItem :=
			pathStyle.Copy().
				Width(m.viewport.Width - lipgloss.Width(fileSize) - lipgloss.Width(box)).
				Render(
					truncate.StringWithTail(
						file.Path, uint(m.viewport.Width-lipgloss.Width(fileSize)-lipgloss.Width(box)), "..."),
				)

		row := lipgloss.JoinHorizontal(lipgloss.Top, box, directoryItem, fileSize)

		s += fmt.Sprintln(row)
	}
	return s
}

func (m model) headerView() string {
	var s string
	if !m.searchFinished {
		dirs := strings.Join(m.searchDirs, ", ")
		s = fmt.Sprintf("%s Scanning directories %s", m.spinner.View(), dirs)
	} else {
		s = "Searched all directories"
	}
	return titleStyle.Render(s)
}

func (m model) footerView() string {
	return footerStyle.Render("Press q to quit")
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
