package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initailize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Hello World"
	p.SetRect(0, 0, 25, 5)
	p.Border = false
	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}

// func updateList(g *gocui.Gui) {
// 	ch := make(chan string)
// 	go Traverse(ch)
// 	for path := range ch {
// 		time.Sleep(time.Millisecond * 100)
// 		g.Update(func(g *gocui.Gui) error {
// 			v, err := g.View("list")
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Fprintln(v, path)
// 			return nil
// 		})
// 	}
// }
