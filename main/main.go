package main

import (
	"github.com/terrywh/keytracker/server"
	"github.com/terrywh/keytracker/data"
	"github.com/terrywh/keytracker/logger"
	"time"
	"os"
	"flag"
)
var svr *server.Server
var dds  data.DataStore

func main() {
	flag.Parse()
	if AppHelp {
		flag.PrintDefaults()
		return
	}
	var err error
	err = logger.Init(AppLogger, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
	logger.Info("engine", ServerEngine, "init ...")
	dds, err = data.New(ServerEngine, AppPath + "/var")
	if err != nil {
		panic(err)
	}
	svr = server.New()
	initHandle()
	initRoutes()

	go svr.ListenAndServe(ServerAddress)
	waitSignal()
	logger.Info("server shutting ...")
	svr.Close()
	time.Sleep(1 * time.Second)
	logger.Info("engine shutdown ...", ServerEngine)
	dds.Close()
	time.Sleep(1 * time.Second)
}
