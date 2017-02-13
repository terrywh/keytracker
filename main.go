package main

import (
	"log"
	"github.com/terrywh/ntracker/tracker"
	"github.com/terrywh/ntracker/config"
	"net"
	"net/http"
	"syscall"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	log.Println("[info] ApiServer started on", config.ApiServerAddr)
	go http.ListenAndServe(config.ApiServerAddr, tracker.GetRESTHandler())


	addr, _   := net.ResolveTCPAddr("tcp", config.NodeServerAddr)
	tl, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal("[error] failed to listen on ", config.NodeServerAddr)
	}
	log.Println("[info] NodeServer started on", config.NodeServerAddr)
	go tracker.GetNodeServer()(tl)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR2)
	config.RotateLogger()
	var s os.Signal
	for {
		s = <- c
		if s == syscall.SIGUSR2 {
			fmt.Fprintln(os.Stderr, "[info] ntracker rotate log file.")
			config.RotateLogger()
		}
	}
}
