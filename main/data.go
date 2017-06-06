package main

import (
	"github.com/terrywh/keytracker/trie"
	"github.com/terrywh/keytracker/server"
	"sync/atomic"
	"encoding/binary"
	"sync"
	"fmt"
	"io"
	"path"
	"time"
)

var dataStore trie.Trie
var dataStoreL *sync.RWMutex

func init() {
	dataStore = trie.NewTrie()
	dataStoreL = &sync.RWMutex{}
}
var keyIncr uint32
func DataKey(key string) string {
	buffer := make([]byte, 6)
	var now uint32
	var inc uint16
	inc = uint16(atomic.AddUint32(&keyIncr, 1))
	now = uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(buffer[0:4], now)
	binary.LittleEndian.PutUint16(buffer[4:6], inc)
	return fmt.Sprintf("%s%02x", key, buffer)
}
func DataKeyFlat(k string) string {
	k = path.Clean(k)
	if k == "/" {
		return ""
	} else {
		return k
	}
}

func DataSet(key string, val interface{}) bool {
	n := dataStore.Get(key)
	if n == nil && val != nil { // 新创建
		dataStore.Create(key).SetValue(val)
		return true // change!
	} else if n == nil && val == nil { // 未变更
		return false
	} else if n != nil && val != nil { // 修改
		return n.SetValue(val)
	} else /*if n!= nil && val == nil */ { // 删除
		dataStore.Remove(key)
		return true
	}
}

func DataDel(key string) bool {
	return dataStore.Remove(key) != nil
}

func DataGet(key string, s io.Writer, y int) {
	n := dataStore.Get(key)
	if n == nil {
		DataWrite(s, key, nil, y)
	}else{
		DataWrite(s, key, n.GetValue(), y)
	}
}

func DataList(key string, s io.Writer, y int, cb func()) {
	n := dataStore.Get(key)
	if n != nil {
		n.Walk(func(c *trie.Node) bool {
			DataWrite(s, key + "/" + c.Key, c.GetValue(), y)
			return true
		})
	}
	if cb != nil {
		cb()
	}
}

func DataWalk(key string, cb func (key string, val interface{}) bool) {
	n := dataStore.Get(key)
	if n != nil {
		n.Walk(func(c *trie.Node) bool {
			return cb(key + "/" + c.Key, c.GetValue())
		})
	}
}

func DataCleanup(s *server.Session) {
	s.WalkElement(func(key string) bool {
		dataStore.Remove(key)
		return true
	})
}

func DataWrite(s io.Writer, key string, val interface{}, y int) {
	if val == nil {
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":null,\"y\":%d}\n", key, y)
		return
	}
	switch val.(type) {
	case float64:
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":%v,\"y\":%d}\n", key, val, y)
	case bool:
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":%t,\"y\":%d}\n", key, val, y)
	default:
		fmt.Fprintf(s, "{\"k\":\"%s\",\"v\":\"%v\",\"y\":%d}\n", key, val, y)
	}
}
