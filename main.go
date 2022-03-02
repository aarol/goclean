package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "goclean"
	app.Usage = "clean up your filesystem"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "exclude, e",
			Value: "node_modules .git",
			Usage: "space separated list of directories to exclude searching into",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "include hidden files",
		},
		cli.BoolFlag{
			Name:  "home",
			Usage: "search from user home instead of current working directory",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "will be written to debug.log",
		},
	}
	app.UsageText = "goclean.exe [options] [directories to search for]"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			return errors.New("no directories to search")
		}
		if c.Bool("debug") {
			f, err := tea.LogToFile("debug.log", "")
			if err != nil {
				return err
			}
			defer f.Close()
		}

		p := tea.NewProgram(initialModel(c), tea.WithAltScreen())
		m, err := p.StartReturningModel()
		if err != nil {
			return err
		}

		return fmt.Errorf("space saved: %s", HumanizeBytes(m.(model).bytesSaved))
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
