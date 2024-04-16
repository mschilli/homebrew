// /////////////////////////////////////////
// imgmod.go - image modifying function
// 2021, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"
)

func imgMod(srcFile string, dstFile string, cb func(image.Image) *image.RGBA) {
	f, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}

	jimg, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}

	dimg := cb(jimg)
	if err != nil {
		log.Fatalf("Modifier failed")
	}

	f, err = os.Create(dstFile)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	err = jpeg.Encode(f, dimg, nil)
	if err != nil {
		log.Fatalf("jpeg.Encode failed: %v", err)
	}
}
