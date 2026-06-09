// ///////////////////////////////////////
// flixer.go - Flixer Main Event Loop
// Mike Schilli, 2026 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	upcoming := flag.Bool("upcoming", false, "show titles containing (upcoming)")
	checkUpcoming := flag.Bool("check-upcoming", false, "check upcoming Netflix titles without starting the UI")
	flag.Parse()

	picker := NewPicker()

	if *checkUpcoming {
		if err := CheckUpcoming(picker.Load()); err != nil {
			panic(err)
		}
		return
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	picks := picker.Load()
	visiblePickIndexes := VisiblePickIndexes(picks, *upcoming)

	lb := widgets.NewList()
	lb.Title = "Movies"
	lb.SelectedRowStyle = ui.NewStyle(ui.ColorBlue)
	lb.TextStyle.Fg = ui.ColorGreen

	info := widgets.NewParagraph()
	info.Title = "Info"
	info.WrapText = true
	info.TextStyle.Fg = ui.ColorBlue

	rate := NewRating()

	grid := ui.NewGrid()
	w, h := ui.TerminalDimensions()
	grid.SetRect(0, 0, w, h)

	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(0.5, lb),
			ui.NewCol(0.5,
				ui.NewRow(0.5, info),
				ui.NewRow(0.5, rate.Widget),
			),
		),
	)

	setInfo := func(pick Pick) {
		service := "No URL"
		u, err := url.Parse(pick.URL)
		if err == nil {
			service = u.Hostname()
		}

		seen := "(not seen)"
		if len(pick.Date) != 0 {
			seen = pick.Date
		}
		info.Text = fmt.Sprintf("Service: %s\nSeen: %s\n",
			service, seen)
	}

	render := func() {
		lb.Rows = []string{}
		for _, idx := range visiblePickIndexes {
			it := picks[idx]
			lb.Rows = append(lb.Rows, it.Title)
		}

		if len(visiblePickIndexes) == 0 {
			info.Text = "No movies"
			rate.Rating = 0
			rate.Update()
			ui.Render(grid)
			return
		}

		if lb.SelectedRow >= len(visiblePickIndexes) {
			lb.SelectedRow = len(visiblePickIndexes) - 1
		}

		it := picks[visiblePickIndexes[lb.SelectedRow]]

		setInfo(it)

		rate.Rating = it.Rating
		rate.Update()

		ui.Render(grid)
	}

	render()

	for e := range ui.PollEvents() {

		switch e.ID {

		case "q", "<C-c>":
			return

		case "e":
			ui.Close()
			runVim(picker.YAMLPath)
			ui.Init()

			picks = picker.Load()
			visiblePickIndexes = VisiblePickIndexes(picks, *upcoming)
			render()

		case "<Resize>":
			pay := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, pay.Width, pay.Height)
			render()

		case "j":
			if lb.SelectedRow < len(visiblePickIndexes)-1 {
				lb.SelectedRow++
				render()
			}

		case "k":
			if lb.SelectedRow > 0 {
				lb.SelectedRow--
				render()
			}

		case "<Enter>":
			if len(visiblePickIndexes) == 0 {
				continue
			}
			idx := visiblePickIndexes[lb.SelectedRow]
			picks[idx].Date = time.Now().Format("2006-01-02")
			picker.Save(picks)
			runFirefox(picks[idx].URL)

		case "w":
			if len(visiblePickIndexes) == 0 {
				continue
			}
			idx := visiblePickIndexes[lb.SelectedRow]
			runFirefox(JustWatchSearchURL(picks[idx].Title))

		case "<MouseLeft>":
			if len(visiblePickIndexes) == 0 {
				continue
			}
			rate.MouseEvent(e.Payload.(ui.Mouse))
			picks[visiblePickIndexes[lb.SelectedRow]].Rating = rate.Rating
			picker.Save(picks)
			render()
		}
	}
}
