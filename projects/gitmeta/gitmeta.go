package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const Version = "0.05"

func main() {
	version := flag.Bool("version", false, "print release version and exit")
	debug := flag.Bool("debug", false, "print verbose debug info")

	flag.Usage = func() {
		fmt.Printf("%s repos.gmf\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), Version)
		return
	}

	config := zapcore.EncoderConfig{
		MessageKey:   "message",
		EncodeTime:   nil,
		EncodeLevel:  nil,
		EncodeCaller: nil,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)
	blog := zap.New(core)
	defer blog.Sync()

	if *debug {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		blog, _ = config.Build()
	}

	log := blog.Sugar()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	cfg := flag.Arg(0)

	gmf := NewGitmeta()
	gmf.Logger = log

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
		err := cloneOrUpdate(log, c, gitDir)
		if err != nil {
			panic(err)
		}
	}
}
