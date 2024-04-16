package main

import (
	"errors"
	"strconv"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

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
	}

	return 0, errors.New("Orientation tag not found")
}
