package main

import (
	"fmt"
	"github.com/terrywh/keytracker/config"
	"github.com/terrywh/keytracker/server"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"net/http/pprof"
)

func main() {
	// 启动一个调试服务器
	pprofSvr := &http.Server{
		Addr:    ":9060",
		Handler: pprof.Handler("heap"),
	}
	go pprofSvr.ListenAndServe()

	go server.ListenAndServe(config.NodeServerAddr, handler, router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	config.RotateLogger()
	var s os.Signal
SIGNAL_WAITING:
	for {
		s = <-c
		switch s {
			case syscall.SIGUSR2:
			fmt.Fprintln(os.Stderr, "[info] keytracker rotate log file.")
			config.RotateLogger()
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			break SIGNAL_WAITING
		}
	}
	fmt.Println("exiting ...")
	dataStoreF.Close()
}
