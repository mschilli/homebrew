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
	app     fyne.App
	cache   *lru.Cache
	mainWin fyne.Window
}

func NewInspector(app fyne.App, mainWin fyne.Window) *inspector {
	insp := inspector{
		app:     app,
		mainWin: mainWin,
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
	entry, ok := i.cache.Get(path.Name())

	var w fyne.Window
	if ok {
		w = entry.(fyne.Window)
	} else {
		var err error
		w, err = i.Load(path)
		if err != nil {
			return err
		}
		i.cache.Add(path.Name(), w)
	}

	w.Hide()
	w.Show()

	return nil
}

func (i *inspector) PreLoad(path fyne.URI) error {
	if !i.cache.Contains(path.Name()) {
		w, err := i.Load(path)
		if err != nil {
			return err
		}
		fyne.Do(func() {
			w.Show()
			i.cache.Add(path.Name(), w)
			w.Hide()
			i.mainWin.Show()
		})
	}
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

	w.SetCloseIntercept(func() {
		w.Hide()
		i.mainWin.Show()
	})

	w.Canvas().SetOnTypedKey(
		func(ev *fyne.KeyEvent) {
			key := string(ev.Name)
			switch key {
			case "Q":
				w.Hide()
				i.mainWin.Show()
			}
		})

	return w, nil
}
