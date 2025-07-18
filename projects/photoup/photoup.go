// /////////////////////////////////////////
// photoup.go - Index files for photo dirs
// 2024, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
)

const Version = "0.01"

func main() {
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *version {
	    fmt.Printf("%s\n", Version)
	    return
	}

	dir, err := photoDir()
	if err != nil {
		panic(err)
	}

	os.MkdirAll(dir, 0755)

	names := []string{}

	tmpl := NewTmpl()
	err = tmpl.Init()
	if err != nil {
		panic(err)
	}

	for _, file := range flag.Args() {
		name, err := processFile(file, dir)
		if err != nil {
			panic(err)
		}

		names = append(names, name)
	}

	idx, err := os.OpenFile(path.Join(dir, "index.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	tpl := NewTmpl()
	if err = tpl.Init(); err != nil {
		panic(err)
	}
	for _, name := range names {
		tpl.AddPhoto(path.Base(name))
	}
	tpl.RenderPage(idx, "photos.html")

	fmt.Printf("%s ready\n", dir)
}

func processFile(fpath string, dir string) (string, error) {
	fh, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	savePath := path.Join(dir, sanitizeFileName(fpath))
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, fh); err != nil {
		return "", err
	}

	if err = scaleJPG(savePath); err != nil {
		return "", err
	}

	return fpath, nil
}

func sanitizeFileName(fileName string) string {
	return path.Base(fileName)
}
