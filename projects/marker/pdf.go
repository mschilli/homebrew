package main

import (
    "fmt"
    "image"
    "bytes"
    "image/png"
	fitz "github.com/gen2brain/go-fitz"
	"github.com/jung-kurt/gofpdf"
)

func openPDF(pdfPath string, page int) (image.Image, error) {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer doc.Close()

	if page < 0 || page >= doc.NumPage() {
		return nil, fmt.Errorf("page out of range")
	}

	img, err := doc.Image(page)
	if err != nil {
		return nil, fmt.Errorf("render page: %w", err)
	}

	return img, nil
}

func writePDF(img image.Image, pdfPath string) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	wPx := img.Bounds().Dx()
	hPx := img.Bounds().Dy()

	wPt := float64(wPx)
	hPt := float64(hPx)

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "pt",
		Size: gofpdf.SizeType{
			Wd: wPt,
			Ht: hPt,
		},
	})

	pdf.AddPage()

	opt := gofpdf.ImageOptions{
		ImageType: "PNG",
	}

	pdf.RegisterImageOptionsReader("img", opt, &buf)

	pdf.ImageOptions("img", 0, 0, wPt, hPt, false, opt, 0, "")

	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		return fmt.Errorf("write pdf: %w", err)
	}

	return nil
}
