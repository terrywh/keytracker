package main

import (
	"os"
	"os/signal"
	"syscall"
)

func waitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<- c
}