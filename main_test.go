package main

import (
	"testing"

	"github.com/aarol/goclean/fs"
	tea "github.com/charmbracelet/bubbletea"
)

// First WindowSizeMsg
func TestWindowInitialSize(t *testing.T) {
	m := model{}
	n, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
	if n.(model).height != 100 {
		t.Fatalf(`model height is %d, want 100`, n.(model).height)
	}
	if cmd != nil {
		t.Fatalf(`model cmd is %v, want nil`, cmd)
	}
}

func TestWindowResize(t *testing.T) {
	m := model{}
	m.viewportReady = true

	n, cmd := m.Update(tea.WindowSizeMsg{Width: 200, Height: 200})

	if n.(model).height != 200 {
		t.Fatalf(`model height is %d, want 200`, n.(model).height)
	}
	if n.(model).viewport.Width != 200 {
		t.Fatalf(`model viewport width is %d, want 200`, n.(model).viewport.Width)
	}
	if cmd != nil {
		t.Fatalf(`model cmd is %v, want nil`, cmd)
	}
}

func TestDirectoryUpdate(t *testing.T) {
	m := model{}
	n, cmd := m.Update(fs.DirEntry{})
	if len(n.(model).directories) != 1 {
		t.Fatalf(`directories length is %d, want 1`, len(n.(model).directories))
	}
	if cmd == nil {
		t.Fatalf(`model cmd is %v, want waitForDirectory`, cmd)
	}
}
