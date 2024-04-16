package main

import (
	"errors"
	"io/ioutil"

	//"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

func exifExtractSegmentList(filename string) (*jpegstructure.SegmentList, error) {
	var sl *jpegstructure.SegmentList

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return sl, err
	}
	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return sl, err
	}

	if intfc == nil {
		return sl, errors.New("no exif")
	}

	sl = intfc.(*jpegstructure.SegmentList)

	return sl, nil
}

func exifReadOrientation(sl *jpegstructure.SegmentList) (uint16, error) {
	rootIfd, _, err := sl.Exif()
	if err != nil {
		return 0, errors.New("exif parse error")
	}

	otag, err := rootIfd.FindTagWithName("Orientation")
	orientation, err := otag[0].Value()

	if err != nil {
		return 0, errors.New("exif orientation value parse error")
	}

	return orientation.([]uint16)[0], nil
}
