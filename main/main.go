package main

import (
	"github.com/terrywh/keytracker/server"
	"github.com/terrywh/keytracker/data"
	"github.com/terrywh/keytracker/logger"
	"time"
)
var svr *server.Server
var dds  data.DataStore

func main() {
	// logger.Init("debug", os.Stdout, os.Stderr);
	var err error
	dds, err = data.New("bolt", AppPath + "/var")
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
