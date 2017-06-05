package main

import (
	"container/list"
	"sync"
	"github.com/terrywh/keytracker/server"
	"strings"
	"log"
)
var watchers map[string]*list.List
var watchersL *sync.RWMutex

func init() {
	watchers = make(map[string]*list.List)
	watchersL = &sync.RWMutex{}
}

func WatcherAppend(key string, s *server.Session) {
	watchersL.Lock()
	defer watchersL.Unlock()

	watcher,ok := watchers[key]
	if !ok {
		watcher = list.New()
		watchers[key] = watcher
	}
	watcher.PushBack(s)

	log.Println("[info] watcher append:", key)
}

func WatcherRemove(key string, s *server.Session) {
	watchersL.Lock()
	defer watchersL.Unlock()
	watcher,ok := watchers[key]
	if !ok {
		return
	}
	for e:=watcher.Front(); e!=nil; e=e.Next() {
		if e.Value.(*server.Session) == s {
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

func WatcherNotify(key string, val interface{}, top string) {
	// /a/b/c => /a/b
	p := strings.LastIndexByte(top, byte('/'))
	if p == -1 {
		return
	}
	top = top[0:p]
	watchersL.RLock()
	defer watchersL.RUnlock()
	watcher,ok := watchers[top]
	if !ok {
		return
	}
	for e:=watcher.Front(); e!=nil; e=e.Next() {
		DataWrite(e.Value.(*server.Session), key, val, 2)
	}
	// 递归通知
	WatcherNotify(key, val, top)
}
