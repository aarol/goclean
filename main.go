package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

func main() {

	if len(os.Args) < 2 {
		log.Panicln("Please provide a directory to search")
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	go updateList(g)

	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("directories", 0, 0, maxX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Searching directories:", strings.Join(os.Args[1:], ", "))
	}
	if v, err := g.SetView("list", 0, 2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Autoscroll = true
	}
	g.SetCurrentView("list")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateList(g *gocui.Gui) {
	ch := make(chan string)
	go Traverse(ch)
	for path := range ch {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("list")
			if err != nil {
				return err
			}
			fmt.Fprintln(v, path)
			return nil
		})
	}
}

// TODO: Fix cursor overflow at bottom
func moveHighlight(v *gocui.View, dy int) error {
	if v != nil {
		// Get the size and position of the view.
		_, y := v.Size()
		ox, oy := v.Origin()

		bottomY := strings.Count(v.ViewBuffer(), "\n") - y - 1
		// If we're at the bottom...
		if oy+dy > bottomY {
			// Set autoscroll to normal again.
			v.Autoscroll = true

			// Move cursor towards the bottom.
			// Only if cursor is not at the bottom already.
			_, cy := v.Cursor()
			v.SetCursor(0, cy+dy)
		} else {
			// Set autoscroll to false and scroll.
			v.Autoscroll = false
			v.SetOrigin(ox, oy+dy)
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
	}
	return nil
}

func initKeybindings(g *gocui.Gui) error {
	// CTRL-C
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	// UP arrow
	if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			moveHighlight(v, -1)
			return nil
		}); err != nil {
		return err
	}
	// DOWN arrow
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			moveHighlight(v, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}
