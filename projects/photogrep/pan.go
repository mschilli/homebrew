// ///////////////////////////////////////
// pan.go - scrolled image panning
// Mike Schilli, 2024 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type Pan struct {
	widget.BaseWidget
	image        *canvas.Image
	scroll       *container.Scroll
	scrollOffset fyne.Position
}

func NewPan(img *canvas.Image) *Pan {
	di := &Pan{image: img}
	di.ExtendBaseWidget(di)

	scrollContent := container.NewMax(di)
	scroll := container.NewScroll(scrollContent)
	di.scroll = scroll

	return di
}

func (di *Pan) CreateRenderer() fyne.WidgetRenderer {
	return &panRenderer{di}
}

func (di *Pan) Dragged(e *fyne.DragEvent) {
	di.scroll.Offset = di.scroll.Offset.SubtractXY(e.Dragged.DX, e.Dragged.DY)
	di.scroll.Refresh()
}

func (di *Pan) Center() {
	w, h := di.image.Size().Width, di.image.Size().Height
	di.scroll.Offset = fyne.NewPos(float32(w)/2.0, float32(h)/2.0)
	di.scroll.Refresh()
}

func (di *Pan) DragEnd() {
}

type panRenderer struct {
	di *Pan
}

func (r *panRenderer) Layout(size fyne.Size) {
	r.di.image.Resize(size)
}

func (r *panRenderer) MinSize() fyne.Size {
	return r.di.image.MinSize()
}

func (r *panRenderer) Refresh() {
	canvas.Refresh(r.di.image)
}

func (r *panRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *panRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.di.image}
}

func (r *panRenderer) Destroy() {}
