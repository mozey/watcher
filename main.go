package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mozey/logutil"
	"github.com/mozey/watcher/pkg/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// version has a hard-coded default value
// https://github.com/mozey/config/issues/20
// For custom builds, the version can be overwritten with ldflags, see
// "Golang compile environment variable into binary"
// https://stackoverflow.com/a/47665780/639133
var version string = "v0.2.0"

func main() {
	// Exit on signal (ctrl + c)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	logutil.SetupLogger(true)

	// Override SetupLogger log level
	debug := os.Getenv("APP_DEBUG") == "true"
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	out, err := watcher.Main(debug)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(2)
	}

	exitCode := 0
	if out.Cmd == watcher.CmdVersion {
		fmt.Println(version)
		sig <- os.Signal(syscall.SIGINT)

	} else if out.Cmd == watcher.CmdWatch {
		defer (func() {
			_ = out.Watcher.Close()
		})()
		go (func() {
			select {
			case err, ok := <-out.Watcher.Errors:
				if !ok {
					return
				}
				log.Error().Stack().Err(err).Msg("")
				exitCode = 2
				sig <- os.Signal(syscall.SIGINT)
			}
		})()
	}

	// Wait on exit signal
	<-sig
	os.Exit(exitCode)
}
