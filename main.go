package main

import (
	"fmt"
	"io"
	"log"
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
			Usage: "include hidden directories",
		},
		cli.BoolFlag{
			Name:  "home",
			Usage: "search from user home instead of current working directory",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "will be written to debug.log in current working directory",
		},
	}
	app.UsageText = "goclean [options] [directories to search for]"
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) == 0 {
			// if no arguments, show help message
			return cli.ShowAppHelp(c)
		}
		if c.Bool("debug") {
			f, err := tea.LogToFile("debug.log", "")
			if err != nil {
				return err
			}
			defer f.Close()
		} else {
			log.SetOutput(io.Discard)
		}

		p := tea.NewProgram(initialModel(c), tea.WithAltScreen())
		m, err := p.StartReturningModel()
		if err != nil {
			return err
		}
		// Print out bytes saved when exiting
		savedStr := fmt.Sprintf("Space saved: %s", HumanizeBytes(m.(model).bytesSaved))
		return cli.NewExitError(savedStr, 0)
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
