package main

import (
	"os"
	"path/filepath"
	"strings"
	"flag"
)

var (
	AppVersion    string = "alpha"
	AppPath       string
	AppLogger     string
	AppHelp       bool
	ServerAddress string
	ServerEngine  string = "bolt"
)

func init() {
	if AppPath == "" {
		AppPath, _ = filepath.Abs(os.Args[0])
		AppPath = filepath.Dir(filepath.Dir(AppPath))
		if strings.HasPrefix(AppPath, "/tmp/") {
			AppPath, _ = os.Getwd();
		}
	}
	flag.BoolVar(&AppHelp, "help", false, "Print this help message")
	flag.StringVar(&ServerAddress, "listen", ":7472", "keytracker will be listening on this addr/port")
	flag.StringVar(&ServerEngine, "engine", "bolt", "data engine, can be one of 'bolt'/'none'; 'bolt' is recommended, 'none' is only used for test purposes.")
	flag.StringVar(&AppLogger, "logger", "debug", "logger level")
	lv := os.Getenv("LOGGER")
	if lv != "" {
		AppLogger = lv
	}
}
