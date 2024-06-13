// /////////////////////////////////////////
// renamer.go - Rename files in bulk
// 2024, Mike Schilli, m@perlmeister.com
// /////////////////////////////////////////
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

const Version = "0.01"

func usage() {
	tmplString := `Usage: {{.Prg}} 'search/replace' file ...

Flags:
{{.Defaults}}
Examples:

$ touch foo1 foo2 foo3

$ {{.Prg}} --verbose "foo/bar" foo*
foo1 -> bar1
foo2 -> bar2
foo3 -> bar3
$ ls bar*
bar1 bar2 bar3

$ {{.Prg}} --verbose "bar/bar-{seq}-" bar*
bar1 -> bar-0001-1
bar2 -> bar-0002-2
bar3 -> bar-0003-3
$ ls bar*
bar-0001-1 bar-0002-2 bar-0003-3

`
	data := struct {
		Prg      string
		Defaults string
	}{
		Prg:      path.Base(os.Args[0]),
		Defaults: captureFlagDefaults(),
	}

	tmpl, err := template.New("").Parse(tmplString)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stderr, data)
	if err != nil {
		panic(err)
	}

	os.Exit(1)
}

func main() {
	dryrun := flag.Bool("dryrun", false, "dryrun only")
	verbose := flag.Bool("verbose", false, "verbose mode")
	help := flag.Bool("help", false, "help message")
	version := flag.Bool("version", false, "print release version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	if *help {
		usage()
	}

	if *dryrun {
		fmt.Printf("Dryrun mode\n")
	}

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Arguments missing\n\n")
		usage()
	}

	cmd := flag.Args()[0]
	files := flag.Args()[1:]
	modifier, err := mkmodifier(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: Invalid command: %s\n", cmd)
		usage()
	}

	for _, file := range files {
		modfile := modifier(file)
		if file == modfile {
			continue
		}
		if *verbose || *dryrun {
			fmt.Printf("%s -> %s\n", file, modfile)
		}
		if *dryrun {
			continue
		}
		err := os.Rename(file, modfile)
		if err != nil {
			fmt.Printf("Renaming %s -> %s failed: %v\n",
				file, modfile, err)
			break
		}
	}
}

func mkmodifier(cmd string) (func(string) string, error) {
	parts := strings.Split(cmd, "/")
	if len(parts) != 2 {
		return nil, errors.New("Invalid repl command")
	}
	search := parts[0]
	repltmpl := parts[1]
	seq := 1

	var rex *regexp.Regexp

	if len(search) == 0 {
		search = ".*"
	}

	rex = regexp.MustCompile(search)

	modifier := func(org string) string {
		repl := strings.Replace(repltmpl,
			"{seq}", fmt.Sprintf("%04d", seq), -1)
		seq++
		res := rex.ReplaceAllString(org, repl)
		return string(res)
	}

	return modifier, nil
}

func captureFlagDefaults() string {
	buf := new(bytes.Buffer)
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	flag.PrintDefaults()

	w.Close()
	os.Stderr = originalStderr
	buf.ReadFrom(r)

	return buf.String()
}
