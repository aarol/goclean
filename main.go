package main

import (
	"context"
	"fmt"
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
				Hidden: true,
				Name:   "fake-entries",
				Usage:  "show fake entries for testing",
			},
		},
		HideHelp:  true,
		UsageText: "goclean [options] [directories to search for]",
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				// if no arguments, show help message
				return cli.ShowAppHelp(c)
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
