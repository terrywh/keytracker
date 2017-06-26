package data

import (
	"sync/atomic"
	// "os"
	// "io/ioutil"
	// "encoding/binary"
	"path"
	"strings"
	"fmt"
	"encoding/json"
	"bytes"
	"github.com/terrywh/keytracker/server"
)
type dataStoreNone struct {
	dbs  map[string][]byte
	inc  uint32
	// ifs *os.File
}
func newDSNone(path string) (*dataStoreNone, error) {
	var ds dataStoreNone
	// var err error
	ds.dbs = make(map[string][]byte)
	ds.inc = 0
	// ds.ifs, err = os.OpenFile(path + "/none.db", os.O_RDWR | os.O_CREATE, os.ModeExclusive | 0666)
	// if err != nil {
	// 	return nil, err
	// }else{
	// 	dt, err := ioutil.ReadAll(ds.ifs)
	// 	if err != nil {
	// 		ds.inc = 0
	// 	}else{
	// 		ds.inc = binary.BigEndian.Uint32(dt)
	// 	}
	// }
	return &ds, nil
}
func (ds *dataStoreNone) Key(k string, suffix bool) string {
	k = path.Clean(k)
	if strings.HasSuffix(k, "/") {
		k = k[0:len(k)-1]
	}
	if !strings.HasPrefix(k, "/") {
		k = "/" + k
	}
	if suffix {
		sq := atomic.AddUint32(&ds.inc, 1)
		// ds.ifs.Seek(0, 0)
		// binary.Write(ds.ifs, binary.BigEndian, sq)
		k = fmt.Sprintf("%s%08x", k, uint32(sq))
	}
	return k
}
func (ds *dataStoreNone) Set(k string, v interface{}) bool {
	vb, ok := ds.dbs[k]
	vc, _ := json.Marshal(v)
	if ok && bytes.Equal(vb, vc) {
		return false
	}
	ds.dbs[k] = vc
	return true
}
func (ds *dataStoreNone) Get(k string) interface{} {
	vb, ok := ds.dbs[k]
	if !ok {
		return nil
	}
	var v server.Request
	json.Unmarshal(vb, &v)
	return v
}
func (ds *dataStoreNone) Del(k string) bool {
	_, ok := ds.dbs[k]
	if !ok {
		return false
	}
	delete(ds.dbs, k)
	return true
}
func (ds *dataStoreNone) List(k string, cb func(key string, val interface{}) bool, r bool) {
	kc := strings.Count(k, "/")
	for ki, vi := range ds.dbs {
		var v server.Request
		json.Unmarshal(vi, &v)
		if v.K == k {
			continue
		} else if !strings.HasPrefix(v.K, k) || (!r && strings.Count(v.K, "/") != kc + 1) || !cb(ki, v) {
			break
		}
	}
}
func (ds *dataStoreNone) Close() error {
	// ds.ifs.Close()
	return nil
}
