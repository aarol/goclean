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
	visibleMax   int
	visibleStart int
	index        int
	init         bool
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
		m.visibleMax = msg.Height
		m.initData()
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			m.moveDown()
		case "up", "k":
			m.moveUp()
		}
	case fs.DirEntry:
		m.addData(msg)
	}
	return m, nil
}

func (m *Model) moveDown() {
	if m.index < m.visibleMax {
		m.index++
	}
	if m.index >= len(m.visibleData) {
		m.shiftVisible(1)
	}
	log.Println(m.index)
}

func (m *Model) shiftVisible(d int) {
	log.Println(m.visibleStart, m.visibleMax, d, len(m.Data))
	if m.visibleStart+m.visibleMax+d < len(m.Data) && m.visibleStart+d > 0 {
		m.visibleStart += d
		m.visibleData = m.Data[m.visibleStart : m.visibleStart+m.visibleMax]
	}
}

func (m *Model) moveUp() {
	if m.index > 0 {
		m.index--
	}
	if m.index < m.visibleStart {
		m.shiftVisible(-1)
	}
	log.Println(m.index)
}

func (m *Model) addData(d fs.DirEntry) {

	if len(m.visibleData) < m.visibleMax {
		m.visibleData = append(m.visibleData, d)
	}
	m.Data = append(m.Data, d)
}

func (m *Model) initData() {
	if m.visibleMax < len(m.Data) {
		m.visibleData = m.Data[:m.visibleMax]
	} else {
		m.visibleData = m.Data[:len(m.Data)]
	}
	m.init = true
}
