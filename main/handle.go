package main

import (
	"github.com/terrywh/keytracker/server"
	"github.com/terrywh/keytracker/logger"
	"sync/atomic"
	"fmt"
	"io"
	"sync"
)
var sessions int32
var mapElement map[*server.Session]map[string]bool
var mapWatcher map[*server.Session]map[string]int
var keySession map[string]map[*server.Session]int
var mapGuard  *sync.Mutex

func initHandle() {
	svr.OnStart = StartHandler
	svr.OnRequest = RequestHandler
	svr.OnClose = CloseHandler
	sessions = 0
	mapElement = make(map[*server.Session]map[string]bool)
	mapWatcher = make(map[*server.Session]map[string]int)
	keySession = make(map[string]map[*server.Session]int)
	mapGuard = &sync.Mutex{}
}
func writeTo(s io.Writer, k string, v interface{}, y int) {
	if v == nil {
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":null,\"y\":%d}\n", k, y)
		return
	}
	switch v.(type) {
	case float64:
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":%v,\"y\":%d}\n", k, v, y)
	case bool:
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":%t,\"y\":%d}\n", k, v, y)
	default: // 剩余一律按文本处理
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":\"%v\",\"y\":%d}\n", k, v, y)
	}
}
func StartHandler(s *server.Session) {
	logger.Info("session start", s)
	atomic.AddInt32(&sessions, 1)
	mapGuard.Lock()
	mapElement[s] = make(map[string]bool)
	mapWatcher[s] = make(map[string]int)
	mapGuard.Unlock()
}

func RequestHandler(s *server.Session, r *server.Request) {
	if r.X >= 512 { // 读取数据
		r.K  = dds.Key(r.K, false)
		if (r.X & 0x02) > 0 { // 直接子集
			dds.List(r.K, func(k string, v interface{}) bool {
				writeTo(s, k, v.(server.Request).V, 2)
				return true
			}, true)
		}else if (r.X & 0x01) > 0 { // 循环子集
			dds.List(r.K, func(k string, v interface{}) bool {
				writeTo(s, k, v.(server.Request).V, 2)
				return true
			}, false)
		}else { // 单条
			req := dds.Get(r.K).(server.Request)
			writeTo(s, req.K, req.V, 0)
		}
	}else if r.X >= 256 { // 监控数据
		mapGuard.Lock()
		r.K = dds.Key(r.K, false)
		if r.V == nil { // 监控 删除
			watcherDel(s, r)
		} else if (r.X & 0x02) > 0 { // 监控 循环子集
			dds.List(r.K, func(k string, v interface{}) bool {
				writeTo(s, k, v.(server.Request).V, 2)
				return true
			}, true)
			watcherAdd(s, r)
		} else if (r.X & 0x01) > 0x00 { // 监控 当前项
			req := dds.Get(r.K).(server.Request)
			writeTo(s, req.K, req.V, 2)
			watcherAdd(s, r)
		} else { // 监控 直接子集
			dds.List(r.K, func(k string, v interface{}) bool {
				writeTo(s, k, v.(server.Request).V, 2)
				return true
			}, false)
			watcherAdd(s, r)
		}
		mapGuard.Unlock()
	}else if r.X >= 0 { // 设置数据
		var changed bool
		if r.V == nil || r.X == 0x00 {
			r.K = dds.Key(r.K, false)
			changed = dds.Del(r.K)
			fmt.Println("del", r.K, changed)
		} else if (r.X & 0x02) > 0 { // 永久
			r.K = dds.Key(r.K, (r.X & 0x04) > 0)
			r.X = r.X & 0x06
			changed = dds.Set(r.K, r)
			mapGuard.Lock()
			delete(mapElement[s], r.K)
			mapGuard.Unlock()
		} else if (r.X & 0x01) > 0 { // 临时
			r.K = dds.Key(r.K, (r.X & 0x04) > 0)
			r.X = r.X & 0x05
			changed = dds.Set(r.K, r)
			mapGuard.Lock()
			mapElement[s][r.K] = true
			mapGuard.Unlock()
		}
		if (r.X & 0x04) > 0 { // 后缀响应
			writeTo(s, r.K, r.V, 1)
		}
		if changed { // 变更通知
			changeNotify(r.K, r.V, r.K, 0)
		}
	}
}
func CloseHandler(s *server.Session) {
	logger.Info("session closed", s)
	atomic.AddInt32(&sessions, -1)
	mapGuard.Lock()
	// 临时标记数据清理
	elementClr(s)
	// 监控清理
	watcherClr(s)
	mapGuard.Unlock()
}
func elementClr(s *server.Session) {
	for k, _ := range mapElement[s] {
		if dds.Del(k) {
			changeNotify(k, nil, k, 0)
		}
	}
	delete(mapElement, s)
}
func watcherAdd(s *server.Session, r *server.Request) {
	// 正向监控标记
	mapWatcher[s][r.K] = r.X
	// 反向监控标记
	if _, ok := keySession[r.K]; !ok {
		keySession[r.K] = make(map[*server.Session]int)
	}
	keySession[r.K][s] = r.X
}
func watcherDel(s *server.Session, r *server.Request) {
	delete(mapWatcher[s], r.K)
	if _, ok := keySession[r.K]; ok {
		delete(keySession[r.K], s)
	}
}
func watcherClr(s *server.Session) {
	for k, _ := range mapWatcher[s] {
		if mapSession, ok := keySession[k]; ok {
			delete(mapSession, s)
			if len(mapSession) == 0 {
				delete(keySession, k)
			}
		}
	}
	delete(mapWatcher, s)
}

func changeNotify(o string, v interface{}, k string, l int) {

	// 当前 KEY 的监控
	if mapSession, ok := keySession[o]; ok {
		for s, x := range mapSession {
			fmt.Printf("N: %v %v %v %v\n", k, o, x, l)
			if (x & 0x02) > 0 && l > 0 || (x & 0x01) > 0 && l == 0 || x == 256 && l == 1 {
				writeTo(s, k, v, 2)
			}
		}
	}
	// 处理为上级 KEY
	if o = changeParent(o); o == "" {
		return
	}
	fmt.Println("K:",k, "O:", o)
	changeNotify(o, v, k, l+1)
}
func changeParent(k string) string {
	if k == "/" {
		return ""
	}
	for i:=len(k)-2; i>-1; i-- {
		if k[i] == '/' {
			if i == 0 {
				return "/"
			}
			return k[0:i]
		}
	}
	return ""
}
