// ///////////////////////////////////////
// pdftki.go - pdftk convenience wrapper
// Mike Schilli, 2020 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
)

const Version = "0.01"

func main() {
	var edit = flag.Bool("e", false, "Pop up an editor")
	var version = flag.Bool("v", false, "Show version")
	flag.Parse()
	pdftkArgs := pdftkArgs(flag.Args())

	if *edit {
		editCmd(&pdftkArgs)
	}

	var out bytes.Buffer
	cmd := exec.Command(pdftkArgs[0],
		pdftkArgs[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OK: [%s]\n", out.String())
}
