// /////////////////////////////////////////
// rotate-180 - Rotate an image 180 degrees
// 2021, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"image"
)

func rot180(jimg image.Image) *image.RGBA {
	bounds := jimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	dimg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dimg.Set(width-1-x, height-1-y, jimg.At(x, y))
		}
	}

	return dimg
}
