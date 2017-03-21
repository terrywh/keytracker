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
	s.AddTag(Tag{key, true})

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
	s.WalkTag(func(tag interface{}) bool {
		_tag := tag.(Tag)
		if _tag.IsWatcher {
			WatcherRemove(_tag.Key, s)
		}
		return true
	})
}

func WatcherNotify(key string, val interface{}) {
	// /a/b/c => /a/b
	p := strings.LastIndexByte(key, byte('/'))
	if p == -1 {
		return
	}
	top := key[0:p]
	watchersL.RLock()
	defer watchersL.RUnlock()
	watcher,ok := watchers[top]
	if !ok {
		return
	}
	for e:=watcher.Front(); e!=nil; e=e.Next() {
		DataWrite(e.Value.(*server.Session), key, val, 2)
	}
}
