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
