package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	const width, height = 100, 100
	bg := color.White
	textColor := color.Black

	for i := 1; i <= 10; i++ {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

		label := fmt.Sprintf("%d", i)
		addLabel(img, 35, 55, label, textColor)

		f, err := os.Create(fmt.Sprintf("%d.jpg", i))
		if err != nil {
			panic(err)
		}
		jpeg.Encode(f, img, nil)
		f.Close()
	}
}

func addLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{col},
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(label)
}
