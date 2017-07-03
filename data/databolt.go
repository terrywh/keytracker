package data

import (
	"github.com/boltdb/bolt"
	"github.com/terrywh/keytracker/server"
	"strings"
	"path"
	"encoding/json"
	"bytes"
	"fmt"
	"time"
)
type dataStoreBolt struct {
	db *bolt.DB
}
var bucketName []byte = []byte("keytracker")
func newDSBolt(path string) (*dataStoreBolt, error) {
	var err error
	var ds dataStoreBolt
	ds.db, err = bolt.Open(path + "/bolt.db", 0644, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	var tx *bolt.Tx
	tx, err = ds.db.Begin(true)
	if err != nil {
		return nil, err
	}
	var bk *bolt.Bucket
	bk, err = tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 清理历史数据中的临时项
	cs := bk.Cursor()
	for _, vb := cs.First(); vb != nil; _, vb = cs.Next() {
		var v server.Request
		json.Unmarshal(vb, &v)
		if (v.X & 0x02) == 0x00 {
			cs.Delete()
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &ds, nil
}
func (ds *dataStoreBolt) Key(k string, suffix bool) string {
	k = path.Clean(k)
	if strings.HasSuffix(k, "/") {
		k = k[0:len(k)-1]
	}
	if !strings.HasPrefix(k, "/") {
		k = "/" + k
	}
	if suffix {
		tx, _ := ds.db.Begin(true)
		bk := tx.Bucket(bucketName)
		sq, _ := bk.NextSequence()
		k = fmt.Sprintf("%s%08x", k, uint32(sq))
		tx.Commit()
	}
	return k
}
func (ds *dataStoreBolt) Set(k string, v interface{}) bool {
	tx, _ := ds.db.Begin(true)
	bk := tx.Bucket(bucketName)

	kb := []byte(k)
	vb := bk.Get(kb)
	vc, _ := json.Marshal(v)
	if bytes.Equal(vb, vc) {
		tx.Commit()
		return false
	}
	bk.Put(kb, vc)
	tx.Commit()
	return true
}
func (ds *dataStoreBolt) Get(k string) interface{} {
	tx, _ := ds.db.Begin(false)
	bk := tx.Bucket(bucketName)

	kb := []byte(k)
	vb := bk.Get(kb)
	var v server.Request
	json.Unmarshal(vb, &v)
	tx.Rollback()
	return v
}
func (ds *dataStoreBolt) Del(k string) bool {
	tx, _ := ds.db.Begin(true)
	bk := tx.Bucket(bucketName)

	kb := []byte(k)
	vb := bk.Get(kb)
	if vb == nil {
		tx.Commit()
		return false
	}
	bk.Delete(kb)
	tx.Commit()
	return true
}
func (ds *dataStoreBolt) List(k string, cb func(key string, val interface{}) bool, r bool) {
	tx, _ := ds.db.Begin(false)
	bk := tx.Bucket(bucketName)
	kc := strings.Count(k, "/")
	kb := []byte(k)
	cs := bk.Cursor()
	for ki, vi := cs.Seek(kb); ki != nil; ki, vi = cs.Next() {
		var v server.Request
		json.Unmarshal(vi, &v)
		if v.K == k {
			continue
		} else if !strings.HasPrefix(v.K, k) {
			break
		} else if strings.Count(v.K, "/") != kc + 1 && !r {
			continue
		} else if !cb(string(ki), v) {
			break
		}
	}
	tx.Rollback()
}
func (ds *dataStoreBolt) Close() error {
	ds.db.Sync()
	return ds.db.Close()
}
