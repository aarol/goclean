package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
)

func main() {
	app := cli.Command{
		Name:  "goclean",
		Usage: "clean up your filesystem",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Value:   []string{"node_modules", ".git"},
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "list all directories, including hidden directories",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "will be written to debug.log in current working directory",
			},
			&cli.BoolFlag{
				Hidden: true,
				Name:   "dry-run",
				Usage:  "show fake data for testing",
			},
		},
		HideHelp:  true,
		UsageText: "goclean [options] [directories to search for]",
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
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
			m, err := p.Run()
			if err != nil {
				return err
			}

			savedStr := fmt.Sprintf("Cleaned: %s", HumanizeBytes(m.(model).bytesSaved))
			return cli.Exit(savedStr, 0)
		},
	}
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
