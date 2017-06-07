package main

import (
	"github.com/terrywh/keytracker/config"
	"github.com/terrywh/keytracker/trie"
	"github.com/terrywh/keytracker/server"
	"sync/atomic"
	"sync"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var dataStore trie.Trie
var dataStoreL *sync.RWMutex
var dataStoreF *os.File
var dataStoreK uint32
func init() {
	dataStore = trie.NewTrie()
	dataStoreL = &sync.RWMutex{}
	var err error
	dataStoreF, err = os.OpenFile(config.AppPath + "/etc/key.inc", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		panic("failed to open key.inc file" + err.Error())
	}
	var tmp []byte
	tmp, err = ioutil.ReadAll(dataStoreF)
	if err != nil {
		panic("failed to read key.inc file" + err.Error())
	}
	var tpk uint64
	tpk, err = strconv.ParseUint(string(tmp), 16, 32)
	if err != nil {
		tpk = 0
	}
	// 防止意外写入同步问题，适当跳过部分数值
	dataStoreK = uint32(tpk + 100)
}

func DataKey(key string) string {
	var inc uint32
	inc = atomic.AddUint32(&dataStoreK, 1)
	dataStoreF.Seek(0, 0)
	fmt.Fprintf(dataStoreF, "%08x", inc)
	// 不做 flush 
	return fmt.Sprintf("%s%08x", key, inc)
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
