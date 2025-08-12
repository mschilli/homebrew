// ///////////////////////////////////////
// clickimg.go - Clickable Image Widget
// Mike Schilli, 2024 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type clickImage struct {
	widget.BaseWidget
	image   *canvas.Image
	cbleft  func()
	cbright func()
}

func newClickImage(img *canvas.Image, cbleft func(), cbright func()) *clickImage {
	ci := &clickImage{}
	ci.ExtendBaseWidget(ci)
	ci.image = img
	ci.cbleft = cbleft
	ci.cbright = cbright
	return ci
}

func (t *clickImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.image)
}

func (t *clickImage) Tapped(_ *fyne.PointEvent) {
	t.cbleft()
}

func (t *clickImage) TappedSecondary(_ *fyne.PointEvent) {
	t.cbright()
}
