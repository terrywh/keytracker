package main

import (
	"os"
	"os/signal"
)

func waitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c
}
