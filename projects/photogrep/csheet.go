// ///////////////////////////////////////
// csheet.go - Contact sheet widget
// Mike Schilli, 2024 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"github.com/disintegration/imaging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	con "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const ThumbSize = float32(200)
const ViewSize = float32(800)

type csheet struct {
	selected map[string]bool
	app      fyne.App
}

func NewCsheet(app fyne.App) *csheet {
	return &csheet{
		selected: map[string]bool{},
		app:      app,
	}
}

func (cs *csheet) MakeGrid(fileNames []string, cache map[string]bool) *con.Scroll {
	grid := con.NewGridWithColumns(3)

	scroll := con.NewScroll(grid)

	go func() {
		for _, fileName := range fileNames {
			pick := cs.newPick(fileName, cache[fileName])
			grid.Add(pick)
			grid.Refresh()
		}
	}()

	return scroll
}

func (cs *csheet) newPick(fileName string, selected bool) *fyne.Container {
	img, err := imaging.Open(fileName, imaging.AutoOrientation(true))
	if err != nil {
		panic(err)
	}
	thumbnail := imaging.Thumbnail(img, int(ThumbSize), int(ThumbSize), imaging.Lanczos)
	image := canvas.NewImageFromImage(thumbnail)

	image.FillMode = canvas.ImageFillContain
	image.SetMinSize(fyne.NewSize(ThumbSize, ThumbSize))

	check := widget.NewCheck("", func(checked bool) {
		if checked {
			cs.selected[fileName] = true
		} else {
			delete(cs.selected, fileName)
		}
	})
	check.SetChecked(selected)

	ci := newClickImage(image, func() {
		// toggle
		check.SetChecked(!check.Checked)
	},
		func() {
			fullView := cs.app.NewWindow("Full Image View")
			img, err := imaging.Open(fileName, imaging.AutoOrientation(true))
			if err != nil {
				panic(err)
			}
			fullImage := canvas.NewImageFromImage(img)
			fullImage.FillMode = canvas.ImageFillOriginal
			inspector := NewPan(fullImage)

			fullView.SetContent(inspector.scroll)
			fullView.Resize(fyne.NewSize(ViewSize, ViewSize))
			inspector.Center()
			fullView.Show()
		})

	image.Refresh()
	return con.NewVBox(ci, check)
}
