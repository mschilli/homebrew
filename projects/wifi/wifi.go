// ///////////////////////////////////////
// wifi.go - term network diagnose tool
// Mike Schilli, 2023 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"github.com/rivo/tview"
	"os"
	"path"
	"strings"
)

var Version = "0.02"

func main() {
	version := flag.Bool("version", false, "print version info")

	flag.Usage = func() {
		fmt.Printf("Usage: %s\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	table.SetBorder(true).SetTitle("Wifi Monitor v" + Version)

	newPlugin(app, table, "Time", clock)
	newPlugin(app, table, "Ping", ping, "www.google.com")
	newPlugin(app, table, "Ping", ping, "8.8.8.8")
	newPlugin(app, table, "Ping", ping, "69.147.88.7")
	newPlugin(app, table, "Ifconfig", nifs)
	newPlugin(app, table, "HTTP", httpGet, "https://youtu.be")

	err := app.SetRoot(table, true).SetFocus(table).Run()
	if err != nil {
		panic(err)
	}
}

func newPlugin(app *tview.Application, table *tview.Table,
	field string, fu func(...string) chan string, arg ...string) {

	if len(arg) > 0 {
		field += " " + strings.Join(arg, " ")
	}

	row := table.GetRowCount()
	table.SetCell(row, 0, tview.NewTableCell(field))

	ch := fu(arg...)

	go func() {
		for {
			select {
			case val := <-ch:
				app.QueueUpdateDraw(func() {
					table.SetCell(row, 1, tview.NewTableCell(val))
				})
			}
		}
	}()
}
