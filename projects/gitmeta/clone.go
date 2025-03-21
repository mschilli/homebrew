package main

import (
	"fmt"
	"github.com/mittwingate/expectre"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

const GitPath = "/usr/bin/git"

func cloneOrUpdate(slog *zap.SugaredLogger, dryrun bool,
	c Cloneable, gitPath string) error {
	fullPath := path.Join(gitPath, c.Dir)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		action := "Cloning"
		if dryrun {
			action += " (dryrun)"
		}
		slog.Infow(action, "repo", c.URL, "dir", fullPath)
		if dryrun {
			return nil
		}
		err := yoyo(slog, GitPath, "clone", c.URL, fullPath)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(fullPath)
	if err != nil {
		return err
	}

	action := "Updating"
	if dryrun {
		action += " (dryrun)"
	}
	slog.Infow(action, "dir", fullPath)
	if dryrun {
		return nil
	}

	err = yoyo(slog, GitPath, "submodule", "foreach", GitPath, "pull")
	if err != nil {
		return err
	}

	err = yoyo(slog, GitPath, "pull")
	if err != nil {
		return err
	}

	return nil
}

const Timeout = 30
const MaxTries = 3

func yoyo(slog *zap.SugaredLogger, args ...string) error {
	slog.Debugw("yoyo", "args", args, "timeout", Timeout, "maxtries", MaxTries)
restart:
	for try := 0; try <= MaxTries; try++ {
		slog.Debugw("Starting", "cmd", args)
		exp := expectre.New()
		err := exp.Spawn(args...)
		if err != nil {
			return err
		}

		var triggered bool

		for {
			select {
			case val := <-exp.Stdout:
				slog.Info(val)
			case val := <-exp.Stderr:
				slog.Info(val)
			case <-time.After(time.Duration(Timeout) * time.Second):
				slog.Infow("Timeout, shutting down", "cmd", args[0])
				triggered = true
				exp.Cancel()
			case <-exp.Released:
				slog.Debugw("Ended", "cmd", args[0])
				if triggered {
					continue restart
				}
				if exp.ExitCode != 0 {
					return fmt.Errorf("%v failed with exit code %d", args, exp.ExitCode)
				}
				break restart
			}
		}
	}

	return nil
}
