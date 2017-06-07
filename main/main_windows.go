package main

import (
	"github.com/terrywh/keytracker/config"
	"github.com/terrywh/keytracker/server"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	go server.ListenAndServe(config.NodeServerAddr, handler, router)
	fmt.Println("wait for signal")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("exiting ...")
	dataStoreF.Close()
}
