package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

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
