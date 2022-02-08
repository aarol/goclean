package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	b  = lipgloss.Color("#bfdeaa")
	kb = lipgloss.Color("#fac984")
	mb = lipgloss.Color("#ff998e")
	gb = lipgloss.Color("#ff7e91")
	tb = lipgloss.Color("#ea45b1")
)

func ColorFromSize(size int64) lipgloss.Color {
	if size < 1024 {
		return b
	} else if size < 1_048_576 {
		return kb
	} else if size < 1_073_741_824 {
		return mb
	} else if size < 1_099_511_627_776 {
		return gb
	}
	return tb
}
