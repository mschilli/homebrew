// ///////////////////////////////////////
// rect.go - Rectangular Utilities
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fyne.io/fyne/v2"
	"image"
	"image/color"
)

type Rect struct {
	From  fyne.Position
	To    fyne.Position
	Zoom  float64
	Color color.NRGBA
}

func NewRect() *Rect {
	return &Rect{
		Color: color.NRGBA{R: 204, G: 255, B: 0, A: 50},
	}
}

func (r *Rect) Bleach() {
	r.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

func (r *Rect) Dims() (fyne.Position, fyne.Size) {
	x := r.From.X
	y := r.From.Y
	w := r.To.X - r.From.X
	h := r.To.Y - r.From.Y

	if r.To.X < r.From.X {
		x = r.To.X
		w = -w
	}

	if r.To.Y < r.From.Y {
		y = r.To.Y
		h = -h
	}

	return fyne.NewPos(x, y), fyne.NewSize(w, h)
}

func (r *Rect) AsImage(zoom float64) image.Rectangle {
	pos, size := r.Dims()

	x := pos.X * float32(zoom)
	y := pos.Y * float32(zoom)
	w := size.Width * float32(zoom)
	h := size.Height * float32(zoom)

	rect := image.Rectangle{
		Min: image.Point{X: int(x), Y: int(y)},
		Max: image.Point{X: int(x) + int(w), Y: int(y) + int(h)},
	}

	return rect
}
