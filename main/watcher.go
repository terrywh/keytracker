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
	r     bool
}
var watchers map[string]*list.List
var watchersL *sync.RWMutex

func init() {
	watchers = make(map[string]*list.List)
}

func WatcherAppend(key string, s *server.Session, r bool) {
	watcher,ok := watchers[key]
	if !ok {
		watcher = list.New()
		watchers[key] = watcher
	}
	watcher.PushBack(WatcherType{s, r})

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

func WatcherNotify(key string, val interface{}, top string, y int) {
	// /a/b/c => /a/bmake
	watcher,ok := watchers[top]
	if ok {
		for e:=watcher.Front(); e!=nil; e=e.Next() {
			wt := e.Value.(WatcherType)
			if y == 0x02 || wt.r { // 顶级 key 或支持递归
				DataWrite(wt.sess, key, val, y)
			}
		}
	}
	// 递归通知
	p := strings.LastIndexByte(top, byte('/'))
	if p > -1 {
		top = top[0:p]
		WatcherNotify(key, val, top, y | 0x04)
	}
}
