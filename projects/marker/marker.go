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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"os"
	"path"
)

const Version = "0.05"

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
	bleach := flag.Bool("bleach", false, "bleach instead of highlight")
	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	if flag.NArg() != 1 {
		usage()
	}

	a := app.NewWithID("com.example.imagehighlighter")
	w := a.NewWindow("Image Highlighter")
	w.Resize(fyne.NewSize(Width, Height))

	ov := NewOverlay()

	if *bleach {
		ov.Bleach()
	}

	imgPath := flag.Args()[0]

	big, img := ov.LoadImage(imgPath)

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
			big, img = ov.LoadImage(imgPath)
			stack.Objects[0] = img
			stack.Refresh()
		}, w).Show()
	})

	saveBtn := widget.NewButton("Save", func() {
		ov.SaveBig(big, imgPath)
	})

	saveAsBtn := widget.NewButton("Save As", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()

			path := writer.URI().Path()
			ov.SaveBig(big, path)
		}, w)
	})

	quitBtn := widget.NewButton("Quit", func() {
		os.Exit(0)
	})

	buttons := container.NewHBox(openBtn, saveBtn, saveAsBtn, quitBtn)
	w.SetContent(container.NewVBox(buttons, stack))
	w.ShowAndRun()
}
