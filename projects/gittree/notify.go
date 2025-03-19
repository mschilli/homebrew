// ///////////////////////////////////////
// notify.go - git repo file sys notifier
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

type Notifier struct {
	watcher *fsnotify.Watcher
	logger  *zap.Logger
}

func NewNotifier(logger *zap.Logger) *Notifier {
	return &Notifier{logger: logger}
}

func (n *Notifier) Wait(minWait time.Duration) {
	var timeoutCh <-chan time.Time

	n.logger.Debug("Waiting ...", zap.Duration("min", minWait))

	// Collect events until minWait timer has elapsed
	for {
		select {
		case ev := <-n.watcher.Events:
			if ev.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove) != 0 {
				n.logger.Debug("Notify write event", zap.String("name", ev.Name))

				if ev.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(ev.Name)
					if err != nil {
						// report but ignore, could be race
						n.logger.Error("Can't stat", zap.String("dir", ev.Name))
						continue
					}

					if info.IsDir() {
						n.logger.Debug("Adding new", zap.String("dir", ev.Name))
						n.watcher.Add(ev.Name)
					}
				}
				if minWait != 0 {
					timeoutCh = time.NewTimer(minWait).C
				}
			}
			continue

		case <-n.watcher.Errors:
		case <-timeoutCh:
			n.logger.Debug("Notify timer expired")
			return
		}
	}
}

func (n *Notifier) Start(rootDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	hashLocal, err := gitHashFileLocal()
	if err != nil {
		return err
	}
	hashRemote, err := gitHashFileRemote()
	if err != nil {
		return err
	}

	n.logger.Debug("Watching", zap.String("hash-local", hashLocal))
	watcher.Add(hashLocal)
	n.logger.Debug("Watching", zap.String("hash-remote", hashRemote))
	watcher.Add(hashRemote)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if filepath.Base(path) == ".git" {
				return filepath.SkipDir // avoid .git loops
			}
			if err := watcher.Add(path); err != nil {
				return err
			}
		}
		return nil
	})

	n.watcher = watcher
	return nil
}
