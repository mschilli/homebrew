// ///////////////////////////////////////
// overlay.go - Drawable canvas
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
	"io"
)

type Overlay struct {
	widget.BaseWidget
	con        *fyne.Container
	markers    []*canvas.Rectangle
	rect       *Rect
	inMotion   bool
	zoom       float64
	CmdPressed bool
}

func NewOverlay() *Overlay {
	over := &Overlay{}
	over.ExtendBaseWidget(over)

	over.con = container.NewWithoutLayout()
	over.rect = NewRect()

	return over
}

func (t *Overlay) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.con)
}

func (t *Overlay) Dragged(e *fyne.DragEvent) {
	if t.inMotion == false {
		// just started
		rect := canvas.NewRectangle(t.rect.Color())
		if t.CmdPressed {
			t.markers = append(t.markers, rect)
		} else {
			    for _, marker := range t.markers {
				t.con.Remove(marker)
			    }
			t.markers = []*canvas.Rectangle{rect}
		}

		t.con.Add(rect)
		t.con.Refresh()

		t.inMotion = true
		t.rect.From = e.Position
		return
	}

	t.rect.To = e.Position

	pos, size := t.rect.Dims()
	t.DrawMarker(pos, size)
}

func (t *Overlay) DragEnd() {
	t.inMotion = false
}

func (t *Overlay) DrawMarker(pos fyne.Position, size fyne.Size) {
	lidx := len(t.markers) - 1
	t.markers[lidx].Resize(size)
	t.markers[lidx].Move(pos)
	t.con.Refresh()
}

func (t *Overlay) SaveBig(big image.Image, path string) error {
	dimg := imaging.Clone(big)

	for _, fyRect := range t.markers {
		rect := NewRect()
		rect.From = fyRect.Position()
		rect.To = fyne.NewPos(
			rect.From.X+fyRect.Size().Width,
			rect.From.Y+fyRect.Size().Height,
		)
		r := rect.AsImage(t.zoom)
		draw.Draw(dimg, r, &image.Uniform{t.rect.Color()}, r.Min, draw.Over)
	}

	err := imaging.Save(dimg, path)
	if err != nil {
		return err
	}

	return nil
}

func (t *Overlay) loadImage(r io.Reader) (image.Image, *canvas.Image) {
	big, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		panic(err)
	}

	shrunk := imaging.Resize(big, Width, 0, imaging.Lanczos)
	t.zoom = float64(big.Bounds().Dx()) / float64(Width)

	img := canvas.NewImageFromImage(shrunk)
	img.FillMode = canvas.ImageFillOriginal
	return big, img
}
