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
	ServerAddress string
)

func init() {
	if AppPath == "" {
		AppPath, _ = filepath.Abs(os.Args[0])
		AppPath = filepath.Dir(filepath.Dir(AppPath))
		if strings.HasPrefix(AppPath, "/tmp/") {
			AppPath, _ = os.Getwd();
		}
	}
	flag.StringVar(&ServerAddress, "listen", ":7472", "keytracker will be listening on this addr/port")
}
