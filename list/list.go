package list

import (
	"fmt"
	"log"

	"github.com/aarol/goclean/fs"
	"github.com/aarol/goclean/util"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Data []fs.DirEntry
	// Obj is the object that was selected, index is the index of the object
	SelectedFunc func(m Model, obj fs.DirEntry, index int)

	visibleData []fs.DirEntry
	// Last index that is visible at a time, or the height
	maxIndex int
	index    int
	init     bool
}

func (m Model) View() string {
	// cursor := termenv.String().Foreground(termenv.ANSIWhite).Bold().String()
	if !m.init {
		return ""
	}
	var data string
	for i, d := range m.visibleData {
		var cursor string
		if i == m.index {
			cursor = "> "
		} else {
			cursor = "  "
		}
		size := util.HumanizeBytes(d.Size)

		data += fmt.Sprintf("\n%s%5s %s", cursor, size, d.Path)
	}
	return data
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.maxIndex = msg.Height
		m.initData()
	case tea.KeyMsg:
		switch msg.String() {

		case "down":
			m.moveDown()
		case "up":
			m.moveUp()
		}
	case fs.DirEntry:
		m.updateData(msg)
	}
	return m, nil
}

func (m *Model) SetHeight(h int) {
	m.maxIndex = h
}

func (m *Model) moveDown() {
	m.index++
}

func (m *Model) moveUp() {
	m.index--
}

func (m *Model) updateData(d fs.DirEntry) {
	m.Data = append(m.Data, d)
	log.Println(m.Data)
	m.visibleData = m.Data
}

func (m *Model) initData() {
	log.Println("initData")
	if m.maxIndex < len(m.Data) {
		m.visibleData = m.Data[:m.maxIndex]
	} else {
		m.visibleData = m.Data[:len(m.Data)]
	}
	m.init = true
}
