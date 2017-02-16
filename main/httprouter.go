package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/terrywh/ntracker/config"
	"runtime"
	"runtime/pprof"
	"encoding/json"
	"github.com/gorilla/websocket"

	"github.com/terrywh/ntracker/server"
)

var router *httprouter.Router
var upgrader websocket.Upgrader


func init() {
	router = httprouter.New()
	router.NotFound = http.FileServer(http.Dir(config.AppPath + "/www"))

	router.GET("/data/:key",  dataGet)
	router.HEAD("/data/:key", dataList)

	router.GET("/session", sessionGet)

	router.GET("/status", statusRouter)
	router.GET("/status/memprof", statusMemProfRouter)


	upgrader = websocket.Upgrader{}
}

func dataGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("{}"))
}
func dataList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("[]"))
}

func sessionGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	s := server.NewSession(server.WrapWebsocket(conn), conn.RemoteAddr().String())

	s.Start(handler)
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

func statusRouter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var status _status
	status.App.Version = config.AppVersion
	status.App.Path    = config.AppPath
	status.App.Env     = config.AppEnv
	status.App.ServerAddr = config.NodeServerAddr
	// TODO 连接数量
	status.App.Sessions  = sessions

	status.Go.Version  = runtime.Version()
	status.Go.Routine  = runtime.NumGoroutine()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	status.Go.HeapAlloc	= ms.HeapAlloc
	status.Go.HeapObjects = ms.HeapObjects

	mm, _ := json.MarshalIndent(status, "", "    ")
	w.Header().Set("Content-Type", "text/json")
	w.Write(mm)
}

func statusMemProfRouter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=ntracker_memory_profile.pprof")

	pprof.WriteHeapProfile(w)
}