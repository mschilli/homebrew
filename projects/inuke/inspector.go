package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/disintegration/imaging"
	"github.com/hashicorp/golang-lru"
	"github.com/mschilli/fyne-loupe"
	"path/filepath"
)

const ViewSize = 1000

type inspector struct {
	app   fyne.App
	cache *lru.Cache
}

func NewInspector(app fyne.App) *inspector {
	insp := inspector{
		app: app,
	}

	// TODO: NewWithEvict
	cache, err := lru.New(10)
	if err != nil {
		panic(err)
	}

	insp.cache = cache

	return &insp
}

func (i *inspector) Show(path fyne.URI) error {
	w, err := i.Load(path)
	if err != nil {
		return err
	}

	w.Show()

	return nil
}

func (i *inspector) Load(path fyne.URI) (fyne.Window, error) {
	w := i.app.NewWindow(filepath.Base(path.Name()))

	img, err := imaging.Open(path.Name(), imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	cimg := canvas.NewImageFromImage(img)

	l := loupe.NewLoupe(cimg)
	w.SetContent(l.Scroll)
	w.Resize(fyne.NewSize(ViewSize, ViewSize))
	l.Center()
	w.Hide()

	return w, nil
}
