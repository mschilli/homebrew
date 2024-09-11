package main

import (
	"flag"
	"fmt"
	"github.com/mschilli/go-murmur"
	"os"
	"path"
	"strings"
	"text/template"
)

const Version = "0.0.1"

func main() {
	fromFile := flag.String("from-file", "", "alternate .murmur file path")
	version := flag.Bool("version", false, "print release version and exit")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [--from-file=file] template-file\n",
			path.Base(os.Args[0]))

		os.Exit(1)
	}

	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
	}

	m := murmur.NewMurmur()

	if len(*fromFile) != 0 {
		m = m.WithFilePath(*fromFile)
	}

	// Go's text/template won't allow variable names containing dashes,
	// add alternates with underscores in the map
	err := m.Read()
	if err != nil {
		panic(err)
	}
	for k, v := range m.Dict {
		dash := "-"
		if strings.Contains(k, dash) {
			newKey := strings.Replace(k, dash, "_", -1)
			m.Dict[newKey] = v
		}
	}

	for _, tmpl := range flag.Args() {
		t, err := template.ParseFiles(tmpl)
		err = t.Execute(os.Stdout, m.Dict)
		if err != nil {
			panic(err)
		}
	}
}
