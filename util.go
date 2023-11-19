package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func HumanizeBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	unit := 0
	for bytes >= 1024 {
		bytes /= 1024
		unit++
	}
	return fmt.Sprintf("%d%s", bytes, units[unit])
}

const (
	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
	tb = gb * 1024
)

func ColorFromSize(size int64) lipgloss.Color {
	if size < kb {
		return lipgloss.Color("#bfdeaa") // byte color
	} else if size < mb {
		return lipgloss.Color("#fac984") // kilobyte color
	} else if size < gb {
		return lipgloss.Color("#ff998e") // megabyte color
	} else if size < tb {
		return lipgloss.Color("#ff7e91") // gigabyte color
	}
	return lipgloss.Color("#ea45b1") // petabyte color
}
