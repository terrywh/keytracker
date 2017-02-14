package tracker

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/terrywh/ntracker/config"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"runtime/pprof"
	"reflect"
	"sync/atomic"
)

var (
	router *httprouter.Router
)

func init() {
	router = httprouter.New()
	router.NotFound = http.FileServer(http.Dir(config.AppPath + "/www"))

	router.GET("/data/:ns/:key", dataRouterItem)
	router.GET("/data/:ns", dataRouterAll)
	router.GET("/data", dataRouterList)

	router.POST("/update/:ns/:key", updateRouter)
	router.GET("/status", statusRouter)
	router.GET("/status/memprof", statusMemProfRouter)
}

func GetRESTHandler() http.Handler {
	return router
}

func dataRouterItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	atomic.AddUint64(&apiStatus.Data, 1)
	w.Header().Set("Content-Type", "text/json")
	found := false
	nlock.RLock() // nodes 的操作需要加锁
	defer nlock.RUnlock()
	for e := nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*_node)
		if n.ns == p[0].Value && n.key == p[1].Value { // 指定的 key 参数
			fmt.Fprint(w, `{"key":"`, n.key, `"`)
			for k, v := range n.data {
				fmt.Fprintf(w, `,"%s":`, k)
				vv, err := json.Marshal(v)
				if err != nil {
					w.Write([]byte("null"))
				}else{
					w.Write(vv)
				}
			}
			w.Write([]byte("}"))
			found = true
		}
	}

	if !found {
		w.WriteHeader(404)
		w.Write([]byte("null"))
	}
}
func dataRouterAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	atomic.AddUint64(&apiStatus.Data, 1)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("["))
	c := 0
	nlock.RLock() // nodes 的操作需要加锁
	defer nlock.RUnlock()
	for e := nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*_node)
		if n.ns == p[0].Value { // 指定的 ns 参数
			if c > 0 { // 首个元素不加 ，逗号
				w.Write([]byte(","))
			}
			c ++
			fmt.Fprint(w, `{"key":"`, n.key, `"`)
			for k, v := range n.data {
				fmt.Fprintf(w, `,"%s":`, k)
				vv, err := json.Marshal(v)
				if err != nil {
					w.Write([]byte("null"))
				}else{
					w.Write(vv)
				}
			}
			w.Write([]byte("}"))
		}
	}
	w.Write([]byte("]"))
}

func dataRouterList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	atomic.AddUint64(&apiStatus.Data, 1)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte("["))
	c := 0
	nlock.RLock() // nodes 的操作需要加锁
	defer nlock.RUnlock()
	for e := nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*_node)
		if c > 0 {
			w.Write([]byte(","))
		}
		c ++
		fmt.Fprintf(w, `{"ns":"%s","key":"%s"}`, n.ns, n.key)
	}
	w.Write([]byte("]"))
}

func updateRouter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	atomic.AddUint64(&apiStatus.Update, 1)
	w.Header().Set("Content-Type", "text/json")

	nlock.RLock() // nodes 的操作需要加锁
	defer nlock.RUnlock()
	for e := nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*_node)
		if n.ns == p[0].Value && n.key == p[1].Value {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("illegal post data - missing"))
				return
			}
			data := make(map[string]interface{})
			err = json.Unmarshal(b, &data)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("illegal post data - json"))
				return
			}
			// 校验赋值的数据类型必须吻合
			for k1,v1 := range n.data {
				v2, ok := data[k1];
				if ok && v2 != nil && reflect.TypeOf(v1).Kind() != reflect.TypeOf(v2).Kind() {
				// if v.(type) != n.data[k].(type) {
					w.WriteHeader(400)
					w.Write([]byte("illegal post data - type"))
					return
				}
			}
			n.wlock.Lock();
			u := 0
			for k,v := range data {
				if v == nil {
					delete(n.data, k)
				}else{
					n.data[k] = v
				}
				u++;
			}
			fmt.Fprintf(w, `{"updated":%d}`, u)
			// 将更新推送给个节点
			n.conn.Write([]byte(`{"action":"data","data":`))
			n.conn.Write(b) // 内部加锁
			n.conn.Write([]byte("}\n"))
			n.wlock.Unlock();
		}
	}
}

type _statusApp struct {
	Version	string
	Path		string
	Env		string
}
type _statusNode struct {
	ServerAddr	string
	Count 		uint
}
type _statusApi struct {
	ServerAddr string
	Data	uint64
	Update	uint64
	Status	uint64
}
type _statusGo struct {
	Version		string
	Routine		int
	HeapAlloc	uint64
	HeapObjects	uint64
}
type _status struct {
	App		_statusApp
	Node	_statusNode
	Api	   *_statusApi
	Go		_statusGo
}

var apiStatus _statusApi

func statusRouter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	atomic.AddUint64(&apiStatus.Status, 1)
	var status _status
	status.App.Version = config.AppVersion
	status.App.Path    = config.AppPath
	status.App.Env     = config.AppEnv
	status.Node.ServerAddr = config.NodeServerAddr
	nlock.RLock() // 这里可以不加锁，但是加上也没什么大影响
	status.Node.Count  = uint(nodes.Len())
	nlock.RUnlock()
	apiStatus.ServerAddr   = config.ApiServerAddr
	status.Api = &apiStatus
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
