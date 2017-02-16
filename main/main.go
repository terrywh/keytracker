package main

import (
	"fmt"
	"github.com/terrywh/keytracker/config"
	"github.com/terrywh/keytracker/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	server.ListenAndServe(config.NodeServerAddr, handler, router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR2)
	config.RotateLogger()
	var s os.Signal
	for {
		s = <-c
		if s == syscall.SIGUSR2 {
			fmt.Fprintln(os.Stderr, "[info] keytracker rotate log file.")
			config.RotateLogger()
		}
	}
}
