package main

import (
	"github.com/terrywh/keytracker/config"
	"github.com/terrywh/keytracker/server"
	"os"
	"os/signal"
)

func main() {

	server.ListenAndServe(config.NodeServerAddr, handler, router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
