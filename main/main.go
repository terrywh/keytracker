package main

import (
	"github.com/terrywh/keytracker/server"
	"github.com/terrywh/keytracker/data"
	"github.com/terrywh/keytracker/logger"
	"time"
	"os"
)
var svr *server.Server
var dds  data.DataStore

func main() {
	var err error
	err = logger.Init(AppLogger, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
	dds, err = data.New(ServerEngine, AppPath + "/var")
	if err != nil {
		panic(err)
	}
	svr = server.New()
	initHandle()
	initRoutes()

	go svr.ListenAndServe(ServerAddress)
	waitSignal()
	logger.Info("exiting ...")
	svr.Close()
	time.Sleep(1 * time.Second)
	dds.Close()
	time.Sleep(1 * time.Second)
}
