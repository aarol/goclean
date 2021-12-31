package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aarol/goclean/fs"
	"github.com/aarol/goclean/list"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli"
)

var (
	footerHeight = 3
	headerHeight = 2
)

type finishedMsg struct{}

type model struct {
	list        list.Model
	spinner     spinner.Model
	finished    bool
	searchDirs  []string
	excludeDirs []string
	directories []string
	sub         chan fs.DirEntry
}

func initialModel(c *cli.Context) model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		directories: []string{},
		sub:         make(chan fs.DirEntry),
		searchDirs:  c.Args(),
		excludeDirs: strings.Split(c.String("exclude"), " "),
		list: list.Model{
			Data: []fs.DirEntry{},
		},
		spinner: s,
	}
}

func listenToSearch(m model) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan fs.DirEntry)

		go fs.Traverse(m.searchDirs, m.excludeDirs, ch)
		for v := range ch {
			m.sub <- v
		}
		return finishedMsg{}
	}
}

func waitForDirectory(sub chan fs.DirEntry) tea.Cmd {
	return func() tea.Msg {
		log.Println("Waiting for directory")
		e := <-sub
		log.Println("Got directory")
		return e
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenToSearch(m),
		waitForDirectory(m.sub),
		spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		log.Println(msg.Height, msg.Width)
		newMsg := tea.WindowSizeMsg{
			Height: msg.Height - footerHeight - headerHeight,
			Width:  msg.Width,
		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(newMsg)
		return m, cmd

	case fs.DirEntry:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, tea.Batch(cmd, waitForDirectory(m.sub))

	case finishedMsg:
		m.finished = true
		return m, nil
	case spinner.TickMsg:
		if !m.finished {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		} else {
			log.Println("Tick after finished")
		}
	}
	return m, nil
}

func (m model) View() string {
	str := ""
	if !m.finished {
		dirs := strings.Join(m.searchDirs, ", ")
		str += fmt.Sprintf("%s Scanning directories %s\n", m.spinner.View(), dirs)
	} else {
		str += "Searched all directories\n"
	}

	str += "\n" + m.list.View()
	return str
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

		f, err := tea.LogToFile("debug.log", "")
		if err != nil {
			return err
		}
		defer f.Close()

		p := tea.NewProgram(initialModel(c), tea.WithAltScreen())
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
