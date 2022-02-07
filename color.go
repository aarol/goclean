package main

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

const (
	weight1 = lipgloss.Color("#ecf772")
	weight2 = lipgloss.Color("#fcbd87")
	weight3 = lipgloss.Color("#ff8c90")
	weight4 = lipgloss.Color("#ff6093")
	weight5 = lipgloss.Color("#ff4d94")
)

func ColorFromSize(size int64) lipgloss.Color {
	if size < 1024 {
		return weight1
	} else if size < int64(math.Pow(1024, 2)) {
		return weight2
	} else if size < int64(math.Pow(1024, 3)) {
		return weight3
	} else if size < int64(math.Pow(1024, 4)) {
		return weight4
	}
	return weight5
}
