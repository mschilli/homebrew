// /////////////////////////////////////////
// autorot.go - Rotate images based on exif
// 2021, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path"
)

const Version = "0.02"

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s jpg-file\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	dryrun := flag.Bool("dryrun", false, "Only show what autorot would do")
	version := flag.Bool("version", false, "Display release version")

	flag.Parse()
	if *version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
	}

	jpgFile := flag.Arg(0)

	f, err := os.Open(jpgFile)
	if err != nil {
		panic(err)
	}

	data, err := exif.Decode(f)
	if err != nil {
		panic(err)
	}

	orient, err := data.Get(exif.Orientation)
	if err != nil {
		fmt.Printf("No orientation header found.\n")
		os.Exit(0)
	}

	val, err := orient.Int(0)
	if err != nil {
		panic(err)
	}

	switch val {
	case 3:
		if *dryrun {
			fmt.Printf("Rotate by 180 degrees\n")
		} else {
			imgMod(jpgFile, jpgFile, rot180)
		}
	case 6:
		if *dryrun {
			fmt.Printf("Rotate by 90 degrees\n")
		} else {
			imgMod(jpgFile, jpgFile, rot90)
		}
	default:
		panic(fmt.Sprintf("Unknown orientation: %d\n", val))
	}
}
