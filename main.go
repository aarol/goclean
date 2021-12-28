package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jroimartin/gocui"
	"github.com/urfave/cli"
)

var Context *cli.Context

func main() {

	app := cli.NewApp()
	app.Name = "Goclean"
	app.Usage = "Clean up your filesystem"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "exclude, e",
			Value: "node_modules .git",
			Usage: "Exclude paths",
		},
	}
	app.Action = func(c *cli.Context) error {
		Context = c
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			return err
		}
		defer g.Close()

		g.SetManagerFunc(layout)

		go updateList(g)

		if err := initKeybindings(g); err != nil {
			return err
		}

		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			return err
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("directories", 0, 0, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Searching directories:", strings.Join(Context.Args(), ", "))
	}
	if v, err := g.SetView("list", 0, 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
	}
	g.SetCurrentView("list")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateList(g *gocui.Gui) {
	ch := make(chan DirEntry)
	go Traverse(ch)
	for entry := range ch {
		// Copy path so that it isn't overwritten
		// Because new path can be received before update executes
		e := entry
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("list")
			if err != nil {
				return err
			}
			size := humanize.Bytes(uint64(e.Size))
			spaces := strings.Repeat(" ", 7-len(size))
			fmt.Fprint(v, SizeColor(e.Size), size, "\033[0m", spaces)
			fmt.Fprintln(v, e.Path)
			return nil
		})
	}
}

func moveHighlight(dy int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil && lineBelow(v, dy) {
			v.MoveCursor(0, dy, false)
		}
		return nil
	}
}

func lineBelow(v *gocui.View, dy int) bool {
	_, cy := v.Cursor()
	line, err := v.Line(cy + dy)
	return err == nil && line != ""
}

func deletePath(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, err := v.Line(cy)
	if err != nil {
		return err
	}
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	// CTRL-C
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	// UP
	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, moveHighlight(-1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", 'k', gocui.ModNone, moveHighlight(-1)); err != nil {
		return err
	}

	// DOWN arrow
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, moveHighlight(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", 'j', gocui.ModNone, moveHighlight(1)); err != nil {
		return err
	}

	// DELETE
	if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, deletePath); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeyDelete, gocui.ModNone, deletePath); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", gocui.KeySpace, gocui.ModNone, deletePath); err != nil {
		return err
	}
	return nil
}
