// ///////////////////////////////////////
// photogrep.go - Photo selector pipe tool
// Mike Schilli, 2024 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const Version = "0.02"

func main() {
	version := flag.Bool("version", false, "Print release version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version)
		return
	}

	var fileNames []string
	if flag.NArg() != 0 {
		fileNames = flag.Args()
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fileNames = append(fileNames, strings.TrimSpace(scanner.Text()))
		}
	}

	app := app.New()
	w := app.NewWindow("photogrep")

	cs := NewCsheet(app)
	cache, err := cacheAsMap()
	if err != nil {
		panic(err)
	}
	grid := cs.MakeGrid(fileNames, cache)
	grid.SetMinSize(fyne.NewSize(800, 600))

	submitButton := widget.NewButton("Submit", func() {
		cache := make(map[string]bool)
		for fileName := range cs.selected {
			cache[fileName] = true
			fmt.Println(fileName)
		}
		err := cacheSave(cache)
		if err != nil {
			panic(err)
		}
		w.Close()
	})

	content := container.NewVBox(grid, submitButton)
	w.SetContent(content)

	w.ShowAndRun()
}
