package util

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	AppVersion string = "alpha"
	AppPath string
)

func init() {
	if AppPath == "" {
		AppPath, _ = filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
		if strings.HasPrefix(AppPath, "/tmp/") {
			AppPath, _ = os.Getwd();
		}
	}
}
