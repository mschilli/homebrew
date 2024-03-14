/////////////////////////////////////////
// image.go - image handling routines
// Mike Schilli, 2021 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"github.com/disintegration/imageorient"
	"github.com/nfnt/resize"
	"image"
	"log"
	"strings"
)

const DspWidth = 1200
const DspHeight = 800

func dispDim(w, h int) (dw, dh int) {
	if w > h {
		// landscape
		return DspWidth, DspHeight
	}

	// portrait
	return DspHeight, DspWidth
}

func isImage(file fyne.URI) bool {
	ext :=
		strings.ToLower(file.Extension())

	if ext != ".jpg" && ext != ".jpeg" {
		return false
	}

	return true
}

func scaleImage(
	img image.Image) image.Image {
	dw, dh := dispDim(img.Bounds().Max.X,
		img.Bounds().Max.Y)
	return resize.Thumbnail(uint(dw),
		uint(dh), img, resize.Lanczos3)
}

func preloadImage(file fyne.URI) {
	if Cache.Contains(file) {
		return
	}
	go func() {
		img, err := loadImage(file)
		if err != nil {
			log.Printf("Unabled to load %s: %v\n", file, err)
			return
		}
		Cache.Add(file, img)
	}()
}

func showImage(
	img *canvas.Image, file fyne.URI) {
	e, ok := Cache.Get(file)
	var nimg *canvas.Image
	if ok {
		nimg = e.(*canvas.Image)
	} else {
		var err error
		nimg, err = loadImage(file)
		if err != nil {
			log.Printf("Unabled to load %s: %v\n", file, err)
			return
		}
		Cache.Add(file, nimg)
	}
	img.Image = nimg.Image

	img.FillMode = canvas.ImageFillOriginal

	img.Refresh()
}

func loadImage(
	file fyne.URI) (*canvas.Image, error) {
	img := canvas.NewImageFromResource(nil)

	read, err :=
		storage.OpenFileFromURI(file)
	if err != nil {
		return nil, err
	}

	defer read.Close()
	raw, _, err := imageorient.Decode(read)
	if err != nil {
		return nil, err
	}

	img.Image = scaleImage(raw)

	return img, nil
}
