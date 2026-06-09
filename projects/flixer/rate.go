/////////////////////////////////////////
// rate.go - Four star rating widget
// Mike Schilli, 2026 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
  ui "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
)

type Rating struct {
  x, y, w, h int
  Widget     *widgets.Paragraph
  Rating     int
}

func NewRating() *Rating {
  p := widgets.NewParagraph()
  p.Title = "Rating"
  p.Text = render(0)

  return &Rating{Rating: 0, Widget: p}
}

func (r *Rating) Update() {
  r.Widget.Text = render(r.Rating)
}

func (r *Rating) MouseEvent(me ui.Mouse) {
  in := r.Widget.Inner
  if me.X >= in.Min.X && me.X < in.Max.X &&
    me.Y >= in.Min.Y && me.Y < in.Max.Y {
    idx := (me.X - in.Min.X) / 2

    if idx >= 0 && idx <= 4 {
      r.Rating = idx
      r.Update()
    }
  }
}

func render(rate int) string {
  s := "□ "
  for i := 0; i < 4; i++ {
    if i < rate {
      s += "[★](fg:green) "
    } else {
      s += "☆ "
    }
  }
  return s
}
