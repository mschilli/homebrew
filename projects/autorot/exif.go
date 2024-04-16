package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

func main() {
	data, err := ioutil.ReadFile("portrait.jpg")
	if err != nil {
		panic(err)
	}
	orientation, err := readOrientation(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Orientation: %d\n", orientation)
}

func readOrientation(data []byte) (uint16, error) {
	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return 0, err
	}

	if intfc == nil {
		return 0, errors.New("no exif")
	}

	sl := intfc.(*jpegstructure.SegmentList)

	_, _, et, err := sl.DumpExif()
	if err != nil {
		if err == exif.ErrNoExif {
			return 0, errors.New("no exif")
		}
	}
	for _, tag := range et {
		if tag.TagName == "Orientation" {
			o, err := strconv.Atoi(tag.FormattedFirst)
			if err != nil {
				return 0, err
			}

			return uint16(o), nil
		}
		//fmt.Printf("%s %s\n", tag.TagName, tag.FormattedFirst)
	}
	return 0, errors.New("no Orientation")
}
