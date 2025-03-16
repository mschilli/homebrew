// ///////////////////////////////////////
// gittree.go - git repo term ui status
// Mike Schilli, 2025 (m@perlmeister.com)
// ///////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

const Version = "0.10"

func main() {
	debug := flag.Bool("debug", false, "enable debug to file")
	version := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	logger := initLogger(*debug)

	rootDir := "."
	if flag.NArg() > 0 {
		rootDir = flag.Arg(0)
	}
	err := os.Chdir(rootDir)
	if err != nil {
		panic(err)
	}

	topDir, err := gitTopDir()
	if err != nil {
		fmt.Println("Error: Cannot find git top dir, are you within a git repository?")
		return
	}

	notify := NewNotifier(logger)
	notify.Start(rootDir)

	app, cmds := ui(topDir)

	go func() {
		for {
			logger.Debug("Checking Status")

			statuses, err := gitStatus()
			if err != nil {
				log.Printf("Status error: %v\n", err)
				break
			}
			logger.Debug("status found", zap.Int("nof_entries", len(statuses)))

			pstatus, err := gitPushStatus()
			if err != nil {
				log.Printf("Push status error: %v\n", err)
				break
			}
			logger.Debug("pstatus is", zap.String("pstatus", pstatus))

			cmds <- Cmd{fs: statuses, pstatus: pstatus}

			minWait := 100 * time.Millisecond
			notify.Wait(minWait)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initLogger(debug bool) *zap.Logger {
	if !debug {
		return zap.NewNop()
	}

	logFile := "/tmp/gittree.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	writeSyncer := zapcore.AddSync(file)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)

	return zap.New(core)
}
