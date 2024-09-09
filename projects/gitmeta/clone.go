package main

import (
	"github.com/mittwingate/expectre"
	"log"
	"os"
	"path"
	"time"
)

func cloneOrUpdate(c Cloneable, gitPath string) error {
	fullPath := path.Join(gitPath, c.Dir)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		err := yoyo("/usr/bin/git", "clone", c.URL, fullPath)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(fullPath)
	if err != nil {
		return err
	}
	err = yoyo("/usr/bin/git", "pull")
	return err
}

const Timeout = 30
const MaxTries = 3

func yoyo(args ...string) error {
restart:
	for try := 0; try <= MaxTries; try++ {
		log.Printf("Starting %s ...\n", args[0])
		exp := expectre.New()
		err := exp.Spawn(args...)
		if err != nil {
			return err
		}

		var triggered bool

		for {
			select {
			case val := <-exp.Stdout:
				log.Println(val)
			case val := <-exp.Stderr:
				log.Println(val)
			case <-time.After(time.Duration(Timeout) * time.Second):
				log.Printf("Timed out. Shutting down ...\n")
				triggered = true
				exp.Cancel()
			case <-exp.Released:
				log.Printf("%s ended.\n", args[0])
				if triggered {
					continue restart
				}
				break restart
			}
		}
	}

	return nil
}
