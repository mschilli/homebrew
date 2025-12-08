// /////////////////////////////////////////
// image.go - Image Processing
// 2024, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func thumbName(fileName string) string {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	return name + "_s" + ext
}

func scaleJPG(inputPath string) error {
	outputPath := thumbName(inputPath)

	img, err := imaging.Open(inputPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	maxLen := 1500
	newWidth, newHeight := rescaleDims(
		img.Bounds().Dx(), img.Bounds().Dy(), maxLen)

	thumbnail := imaging.Thumbnail(img, newWidth, newHeight, imaging.Lanczos)
	err = imaging.Save(thumbnail, outputPath)
	return err
}

func rescaleDims(w, h, max int) (int, int) {
	newW := max
	newH := max

	if w <= max && h <= max {
	    return w, h
	}

	if w > h { // landscape
		newH = h * max / w
	} else {
		newW = w * max / h
	}

	return newW, newH
}
