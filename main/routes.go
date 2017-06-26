package main

import (
	"net/http"
	"github.com/terrywh/keytracker/server"
	"encoding/json"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"runtime"
	"github.com/terrywh/keytracker/logger"
)
func initRoutes() {
	svr.Router.NotFound = http.FileServer(http.Dir(AppPath + "/www"))
	svr.Router.GET("/status", routerStat)

	svr.Router.GET("/read/*key", routerRead)
	svr.Router.POST("/push", routerPush)
	svr.Router.GET("/list/*key", routerList)
}

func routerRead(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	key := dds.Key(p[0].Value, false)
	req := dds.Get(key).(server.Request)
	writeTo(w, key, req.V, 0)
}

func routerPush(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("failed to read json:", err)
		w.Write([]byte("{\"error\":\"failed to read json\"}"))
		return
	}
	var req server.Request
	err = json.Unmarshal(raw, &req)
	if err != nil {
		logger.Warn("failed to decode json:", err)
		w.Write([]byte("{\"error\":\"failed to decode json\"}"))
		return
	}
	req.X = req.X & 0x02
	if (req.X & 0x04) > 0 {
		req.K = dds.Key(req.K, true)
	}else{
		req.K = dds.Key(req.K, false)
	}
	var changed bool
	if req.V == nil {
		changed = dds.Del(req.K)
	}else{
		changed = dds.Set(req.K, req)
	}
	if changed {
		changeNotify(req.K, req.V, req.K, 0)
	}
	if (req.X & 0x04) > 0 {
		writeTo(w, req.K, req.V, 1)
	}else{
		writeTo(w, req.K, req.V, 0)
	}
}

func routerList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")

	w.Write([]byte("["))
	n := 0
	key := dds.Key(p[0].Value, false)
	dds.List(key, func (k string, v interface{}) bool {
		if n != 0 {
			w.Write([]byte(","))
		}
		writeTo(w, k, v.(server.Request).V, 0)
		return true
	}, r.FormValue("r") != "")
	w.Write([]byte("]"))
}

type _statusApp struct {
	Version	string
	Path		string
	Env		string
	ServerAddr string
	Sessions   int32
}
type _statusGo struct {
	Version		string
	Routine		int
	HeapAlloc	uint64
	HeapObjects	uint64
}
type _status struct {
	App		_statusApp
	Go		_statusGo
}
func routerStat(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var status _status
	status.App.Version    = AppVersion
	status.App.Path       = AppPath
	status.App.ServerAddr = ServerAddress
	// 连接数量
	status.App.Sessions   = sessions

	status.Go.Version     = runtime.Version()
	status.Go.Routine     = runtime.NumGoroutine()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	status.Go.HeapAlloc	= ms.HeapAlloc
	status.Go.HeapObjects = ms.HeapObjects

	mm, _ := json.MarshalIndent(status, "", "    ")
	w.Header().Set("Content-Type", "text/json")
	w.Write(mm)
}
