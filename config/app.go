package config

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	AppVersion string = "alpha"
	AppPath    string
	AppEnv     string
)

func init() {
	AppEnv = os.Getenv("GOENV")
	if AppEnv == "" {
		AppEnv = "wuhao"
	}

	if AppPath == "" {
		AppPath, _ = filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
		if strings.HasPrefix(AppPath, "/tmp/") {
			AppPath, _ = os.Getwd();
		}
	}

	initConfigFile()
}
