/////////////////////////////////////////
// trash.go - discarding bad images
// Mike Schilli, 2021 (m@perlmeister.com)
/////////////////////////////////////////
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

func stashDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, "inuke")
}

func toStash(file fyne.URI) {
	toDir := stashDir()
	os.MkdirAll(toDir, 0755)

	fmt.Printf("Stashing\n")
	_, err := copy(file.Name(), filepath.Join(toDir, file.Name()))
	if err != nil {
		panic(err)
	}
}

const TrashDir = "old"

func toTrash(file fyne.URI) {
	err := os.MkdirAll(TrashDir, 0755)
	panicOnErr(err)
	err = os.Rename(file.Name(),
		filepath.Join(TrashDir, file.Name()))
	panicOnErr(err)
}

func fromTrash(file fyne.URI) {
	err := os.Rename(filepath.Join(TrashDir, file.Name()), file.Name())
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dest.Close()
	nBytes, err := io.Copy(dest, source)
	return nBytes, err
}
