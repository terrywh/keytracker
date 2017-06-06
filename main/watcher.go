package main

import (
	"container/list"
	"sync"
	"github.com/terrywh/keytracker/server"
	"strings"
	"log"
)
type WatcherType struct {
	sess *server.Session
	x     int
}
var watchers map[string]*list.List
var watchersL *sync.RWMutex

func init() {
	watchers = make(map[string]*list.List)
}

func WatcherAppend(key string, s *server.Session, x int) {
	watcher,ok := watchers[key]
	if !ok {
		watcher = list.New()
		watchers[key] = watcher
	}
	watcher.PushBack(WatcherType{s, x})

	log.Println("[info] watcher append:", key)
}

func WatcherRemove(key string, s *server.Session) {
	watcher,ok := watchers[key]
	if !ok {
		return
	}
	for e:=watcher.Front(); e!=nil; e=e.Next() {
		if e.Value.(WatcherType).sess == s {
			watcher.Remove(e)
		}
	}
	if watcher.Len() == 0 {
		delete(watchers,key)
	}
	log.Println("[info] watcher removed:", key)
}

func WatcherCleanup(s *server.Session) {
	s.WalkWatcher(func(key string) bool {
		WatcherRemove(key, s)
		return true
	})
}

func WatcherNotify(key string, val interface{}) {
	// /a/b/c => /a/bmake
	watcher,ok := watchers[key]
	if ok {
		for e:=watcher.Front(); e!=nil; e=e.Next() {
			wt := e.Value.(WatcherType)
			if (wt.x & 0x01) != 0x00 { // 单项监控
				DataWrite(wt.sess, key, val, 0x02)
			}
		}
	}
	// 通知上级 key
	p := strings.LastIndexByte(key, byte('/'))
	if p > -1 {
		top := key[0:p]
		watcher,ok := watchers[top]
		if ok {
			for e:=watcher.Front(); e!=nil; e=e.Next() {
				wt := e.Value.(WatcherType)
				if (wt.x & 0x01) == 0x00  { // 子集监控
					DataWrite(wt.sess, key, val, 0x02)
				}
			}
		}
	}
}
