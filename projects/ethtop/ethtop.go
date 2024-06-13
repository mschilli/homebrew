///////////////////////////////////////////
// iftop - Terminal UI displays networks
// 2018, Mike Schilli, m@perlmeister.com
///////////////////////////////////////////
package main

import (
	t "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"time"
)

var listItems = []string{}

func main() {
	err := t.Init()
	if err != nil {
		log.Fatalln("Termui init failed")
	}

	// Cleanup UI on exit
	defer t.Close()

	// Listbox displaying interfaces
	lb := widgets.NewList()
	lb.Title = "Networks"

	// Textbox
	txt := widgets.NewParagraph()
	txt.Text = "Type 'q' to quit."

	resize := func() {
		w, h := t.TerminalDimensions()
		lb.SetRect(0, 0, w, h-3)
		txt.SetRect(0, h-3, w, h)
		t.Render(lb, txt)
	}

	resize()

	uiEvents := t.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Resize>":
				resize()
			case "q", "<C-c>":
				return
			}
		case <-time.After(1 * time.Second):
			lb.Rows = NifsAsStrings()
			t.Render(lb)
		}
	}
}
