package main

import (
	"fmt"
	"github.com/mozey/logutil"
	"github.com/mozey/watcher/pkg/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Exit on signal (ctrl + c)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	logutil.SetupLogger(true)

	// Override SetupLogger log level
	debugKey := "APP_DEBUG"
	debug := os.Getenv(debugKey) // TODO Configurable debugKey
	level := zerolog.InfoLevel
	if debug == "true" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	out, err := watcher.Main()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(2)
	}

	exitCode := 0
	if out.Cmd == watcher.CmdVersion {
		// TODO Add version
		fmt.Println("n/a")
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
