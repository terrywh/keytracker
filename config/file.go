package config

import (
	"os"
	"github.com/BurntSushi/toml"
	"time"
	"log"
	"fmt"
)

var (
	ApiServerAddr  string = ":8360"
	NodeServerAddr string = ":7472"
	NodeKeepAlive  time.Duration = 60 * time.Second

	AppLogger  string
	appLogger *os.File
)

func initConfigFile() {
	var err error
	var ok bool
	cf := AppPath + "/etc/tracker." + AppEnv + ".toml"
	cd := make(map[string]map[string]interface{})
	// 读取 对应配文件，填充上面变量
	_, err = toml.DecodeFile(cf, &cd)
	if err != nil {
		log.Printf("[warning] failed to load config file '%s', using defaults.", cf)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("[warning] failed to parse config file:", r)
		}
	}()

	if _, ok = cd["node"]["server_addr"]; ok {
		NodeServerAddr = cd["node"]["server_addr"].(string)
	}
	if _, ok = cd["api"]["server_addr"]; ok {
		ApiServerAddr = cd["api"]["server_addr"].(string)
	}
	if lf, ok := cd["app"]["logger"]; ok {
		AppLogger = lf.(string)
	}
	RotateLogger()
}

func RotateLogger() {
	if appLogger != nil {
		appLogger.Close()
	}
	if AppLogger != "" {
		appLogger, err := os.OpenFile(AppLogger, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0777)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[error] keytracker failed to rotate log file:", err)
			return
		}
		log.SetOutput(appLogger)
	}else{
		log.SetOutput(os.Stdout)
	}
}
