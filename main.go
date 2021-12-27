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
	l, err := v.Line(cy)
	if err != nil {
		return err
	}
	log.Println(l)
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
	if err := g.SetKeybinding("list", 'j', gocui.ModNone, moveHighlight(-1)); err != nil {
		return err
	}
	// DOWN arrow
	if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, moveHighlight(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("list", 'k', gocui.ModNone, moveHighlight(1)); err != nil {
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
