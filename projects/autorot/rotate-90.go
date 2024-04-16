// /////////////////////////////////////////
// rotate-90 - Rotate an image 90 degrees
// 2021, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"image"
)

func rot90(jimg image.Image) *image.RGBA {
	bounds := jimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bounds.Max.X, bounds.Max.Y = bounds.Max.Y, bounds.Max.X

	dimg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			org := jimg.At(x, y)
			dimg.Set(height-y, x, org)
		}
	}

	return dimg
}
