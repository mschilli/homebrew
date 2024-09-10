package main

import (
	"github.com/mittwingate/expectre"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

const GitPath = "/usr/bin/git"

func cloneOrUpdate(slog *zap.SugaredLogger, c Cloneable, gitPath string) error {
	fullPath := path.Join(gitPath, c.Dir)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		slog.Infow("Cloning", "repo", c.URL, "dir", fullPath)
		err := yoyo(slog, GitPath, "clone", c.URL, fullPath)
		if err != nil {
			return err
		}
	}

	err = os.Chdir(fullPath)
	if err != nil {
		return err
	}
	slog.Infow("Updating", "dir", fullPath)
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
				break restart
			}
		}
	}

	return nil
}
