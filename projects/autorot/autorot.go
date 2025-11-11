package main

import (
	"bytes"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"image"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

const Version = "0.03"

var Log *zap.Logger

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s jpg-file\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	dryrun := flag.Bool("dryrun", false, "Only show what autorot would do")
	version := flag.Bool("version", false, "Display release version")
	debug := flag.Bool("debug", false, "Print verbose debug messages")

	flag.Parse()
	if *version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
	}

	jpgFile := flag.Arg(0)

	if *debug || *dryrun {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}

	Log.Debug("Reading", zap.String("path", jpgFile))

	data, err := os.ReadFile(jpgFile)
	if err != nil {
		Log.Error("Read error", zap.Error(err))
	}

	exif.RegisterParsers(mknote.All...)

	exifData, exifBytes := extractEXIF(data)

	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		Log.Error("Decode error", zap.Error(err))
	}

	if exifData != nil {
		var o int
		img, o, err = fixOrientation(img, exifData)
		if err != nil {
			Log.Error("Exif error", zap.Error(err))
			return
		}
		if o == 1 {
			Log.Error("Orientation header set to 1, no rotation.")
			return
		}
	}

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG); err != nil {
		Log.Error("Encode error", zap.Error(err))
	}

	if exifData != nil {
		exifBytes = resetOrientation(exifBytes)
	}

	output := insertEXIF(buf.Bytes(), exifBytes)

	if *dryrun {
		Log.Debug("Exiting without writing in dryrun mode", zap.String("path", jpgFile))
		return
	}

	Log.Debug("Writing file", zap.String("path", jpgFile))
	if err := os.WriteFile(jpgFile, output, 0644); err != nil {
		Log.Error("Write error", zap.Error(err))
	}
}

func fixOrientation(img image.Image, exifData *exif.Exif) (image.Image, int, error) {
	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		return img, 0, err
	}
	orientation, err := tag.Int(0)
	if err != nil {
		return img, 0, err
	}

	Log.Debug("Found exif tag", zap.Int("orientation", orientation))

	switch orientation {
	case 2:
		Log.Debug("Fix: FlipH")
		img = imaging.FlipH(img)
	case 3:
		Log.Debug("Fix: Rotate180")
		img = imaging.Rotate180(img)
	case 4:
		Log.Debug("Fix: FlipV")
		img = imaging.FlipV(img)
	case 5:
		Log.Debug("Fix: Transpose")
		img = imaging.Transpose(img)
	case 6:
		Log.Debug("Fix: Rotate270")
		img = imaging.Rotate270(img)
	case 7:
		Log.Debug("Fix: Transverse")
		img = imaging.Transverse(img)
	case 8:
		Log.Debug("Fix: Rotate90")
		img = imaging.Rotate90(img)
	}

	return img, orientation, nil
}

func extractEXIF(jpegData []byte) (*exif.Exif, []byte) {
	const (
		markerStart = 0xFF
		app1Marker  = 0xE1
	)
	for i := 0; i < len(jpegData)-1; i++ {
		if jpegData[i] == markerStart && jpegData[i+1] == app1Marker {
			if i+4 >= len(jpegData) {
				break
			}
			segLen := int(jpegData[i+2])<<8 | int(jpegData[i+3])
			if i+4+segLen > len(jpegData) {
				break
			}
			seg := jpegData[i+4 : i+4+segLen]
			if bytes.HasPrefix(seg, []byte("Exif\x00\x00")) {
				ex, err := exif.Decode(bytes.NewReader(seg))
				if err == nil {
					return ex, append([]byte(nil), seg...)
				}
			}
		}
	}
	return nil, nil
}

func insertEXIF(jpegData []byte, exifSegment []byte) []byte {
	if exifSegment == nil {
		return jpegData
	}
	const (
		markerStart = 0xFF
		soiMarker   = 0xD8
		app1Marker  = 0xE1
	)
	if len(jpegData) < 2 || jpegData[0] != markerStart || jpegData[1] != soiMarker {
		return jpegData
	}
	length := len(exifSegment) + 2
	app1Header := []byte{
		markerStart, app1Marker,
		byte(length >> 8), byte(length & 0xFF),
	}
	out := make([]byte, 0, len(jpegData)+len(app1Header)+len(exifSegment))
	out = append(out, jpegData[:2]...) // SOI
	out = append(out, app1Header...)
	out = append(out, exifSegment...)
	out = append(out, jpegData[2:]...)
	return out
}

// resetOrientation overwrites the EXIF orientation tag to "1".
func resetOrientation(exifSegment []byte) []byte {
	// Very naive patch: look for bytes 0x01 0x12 (Orientation tag)
	// and set next two bytes (value) to 0x00 0x01 (1)
	b := make([]byte, len(exifSegment))
	copy(b, exifSegment)
	for i := 0; i < len(b)-4; i++ {
		if b[i] == 0x01 && b[i+1] == 0x12 { // Orientation tag ID
			// We can’t assume endian; safest approach: zero out to normal
			for j := 0; j < 4 && i+j+8 < len(b); j++ {
				b[i+j+8] = 0 // zero bytes before value field
			}
			if i+8 < len(b) {
				b[i+8] = 0x00
			}
			if i+9 < len(b) {
				b[i+9] = 0x01 // set to 1
			}
			break
		}
	}
	return b
}
