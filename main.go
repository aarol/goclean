package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

type model struct {
	searchDirs  []string
	excludeDirs []string
	directories []string
	sub         chan DirEntry
	cursor      int
}

func initialModel(c *cli.Context) model {
	return model{
		directories: []string{},
		sub:         make(chan DirEntry),
		searchDirs:  c.Args(),
		excludeDirs: strings.Split(c.String("exclude"), " "),
	}
}

func listenToSearch(m model) tea.Cmd {
	return func() tea.Msg {
		Traverse(m.searchDirs, m.excludeDirs, m.sub)
		log.Println("Traversing")
		return nil
	}
}

func waitForDirectory(sub chan DirEntry) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenToSearch(m),
		waitForDirectory(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.directories)-1 {
				m.cursor++
			}

		}

	case DirEntry:
		if msg.Path != "" {
			m.directories = append(m.directories, msg.Path)
		}
		return m, waitForDirectory(m.sub)
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Searching for directories: " + strings.Join(m.searchDirs, ", ") + "\n\n"

	// Iterate over our choices
	for i, dir := range m.directories {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, dir)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {

	app := cli.NewApp()
	app.Name = "Goclean"
	app.Usage = "Clean up your filesystem"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "exclude, e",
			Value: "node_modules .git",
			Usage: "Space separated list of directories to exclude searching into",
		},
	}
	app.UsageText = "goclean.exe [options] [directories to search for]"
	app.Action = func(c *cli.Context) error {

		p := tea.NewProgram(initialModel(c))
		if err := p.Start(); err != nil {
			return err
		}

		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
