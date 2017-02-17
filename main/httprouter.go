package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/terrywh/keytracker/config"
	"runtime"
	"runtime/pprof"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/terrywh/keytracker/server"	
	"strconv"
)

var router *httprouter.Router
var upgrader websocket.Upgrader


func init() {
	router = httprouter.New()
	router.NotFound = http.FileServer(http.Dir(config.AppPath + "/www"))

	router.GET("/read/*key", routerRead)
	router.GET("/list/*key", routerList)
	router.GET("/some/:limit/*key", routerSome)

	router.GET("/session", sessionGet)

	router.GET("/status", statusRouter)
	router.GET("/status/memprof", statusMemProfRouter)


	upgrader = websocket.Upgrader{}
}

func routerRead(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")

	DataGet(DataKeyFlat(p[0].Value), w)
}
func routerList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("["))
	n := 0
	DataWalk(DataKeyFlat(p[0].Value), func(key string, val interface{}) bool {
		if n != 0 {
			w.Write([]byte(","))
		}
		DataWrite(w, key, val, 0)
		n++
		return true
	})
	w.Write([]byte("]"))
}
func routerSome(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("["))
	n := 0
	c, _ := strconv.Atoi(p[0].Value)
	DataWalk(DataKeyFlat(p[1].Value), func(key string, val interface{}) bool {
		if n != 0 {
			w.Write([]byte(","))
		}
		DataWrite(w, key, val, 0)
		n++
		if n < c {
			return true
		}
		return false
	})
	w.Write([]byte("]"))
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
	// 连接数量
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
	w.Header().Set("Content-Disposition", "attachment; filename=keytracker_memory_profile.pprof")

	pprof.WriteHeapProfile(w)
}
