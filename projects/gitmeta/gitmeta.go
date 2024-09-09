package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const Version = "0.03"

func main() {
	version := flag.Bool("version", false, "print release version and exit")

	flag.Usage = func() {
		fmt.Printf("%s repos.gmf\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
	}

	cfg := flag.Arg(0)

	gmf := NewGitmeta()

	f, err := os.Open(cfg)
	if err != nil {
		panic(err)
	}

	gmf.AddGMF(f)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	gitDir := filepath.Join(usr.HomeDir, "git")
	err = os.Chdir(gitDir)
	if err != nil {
		panic(err)
	}

	for _, c := range gmf.AllCloneables() {
		fmt.Printf("Cloning %v\n", c)
		err := cloneOrUpdate(c, gitDir)
		if err != nil {
			panic(err)
		}
	}
}
