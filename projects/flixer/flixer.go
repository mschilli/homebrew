/////////////////////////////////////////
// flixer.go - Flixer Main Event Loop
// Mike Schilli, 2026 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
  "fmt"
  "net/url"
  "time"

  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
)

func main() {
  if err := ui.Init(); err != nil {
    panic(err)
  }
  defer ui.Close()

  picker := NewPicker()
  picks := picker.Load()

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
    for _, it := range picks {
      lb.Rows = append(lb.Rows, it.Title)
    }

    it := picks[lb.SelectedRow]

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
      render()

    case "<Resize>":
      pay := e.Payload.(ui.Resize)
      grid.SetRect(0, 0, pay.Width, pay.Height)
      render()

    case "j":
      if lb.SelectedRow < len(picks)-1 {
        lb.SelectedRow++
        render()
      }

    case "k":
      if lb.SelectedRow > 0 {
        lb.SelectedRow--
        render()
      }

    case "<Enter>":
      picks[lb.SelectedRow].Date = time.Now().Format("2006-01-02")
      picker.Save(picks)
      runFirefox(picks[lb.SelectedRow].URL)

    case "<MouseLeft>":
      rate.MouseEvent(e.Payload.(ui.Mouse))
      picks[lb.SelectedRow].Rating = rate.Rating
      picker.Save(picks)
      render()
    }
  }
}
