// ///////////////////////////////////////
// marker.go - image/text highlighter
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image"
	"os"
	"path"
)

const Version = "0.04"

const (
	Width  = 800
	Height = 600
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: "+os.Args[0]+" document.png\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	if len(flag.Args()) != 1 {
		usage()
	}

	a := app.NewWithID("com.example.imagehighlighter")
	w := a.NewWindow("Image Highlighter")
	w.Resize(fyne.NewSize(Width, Height))

	ov := NewOverlay()

	img := &canvas.Image{}
	var big image.Image
	var imgPath string

	if len(os.Args) == 2 {
		imgPath = os.Args[1]
		f, err := os.Open(imgPath)
		if err != nil {
			panic(err)
		}
		big, img = ov.loadImage(f)
	}

	stack := container.NewStack(img, ov)

	if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			if ev.Name == desktop.KeySuperLeft { // Cmd on macOS
				ov.CmdPressed = true
			}
		})
		deskCanvas.SetOnKeyUp(func(ev *fyne.KeyEvent) {
			if ev.Name == desktop.KeySuperLeft {
				ov.CmdPressed = false
			}
		})
	}

	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		key := string(ev.Name)
		switch key {
		case "Q":
			os.Exit(0)
		}
	})

	openBtn := widget.NewButton("Open Image", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			imgPath = reader.URI().Path()
			big, img = ov.loadImage(reader)
			stack.Objects[0] = img
			stack.Refresh()
		}, w).Show()
	})

	saveBtn := widget.NewButton("Save", func() {
		ov.SaveBig(big, imgPath)
	})
	quitBtn := widget.NewButton("Quit", func() {
		os.Exit(0)
	})

	buttons := container.NewHBox(openBtn, saveBtn, quitBtn)
	w.SetContent(container.NewVBox(buttons, stack))
	w.ShowAndRun()
}
