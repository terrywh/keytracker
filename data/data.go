package data

import (
	"errors"
)

var ErrDataEngineNotSupported = errors.New("data engine is not suppored")

type DataStore interface {
	Key(k string, suffix bool) string
	Set(k string, v interface{}) bool
	Get(k string) interface{}
	Del(k string) bool
	List(k string, cb func(key string, val interface{}) bool, r bool)
	Close() error
}

func New(which, path string) (DataStore, error) {
	if which == "bolt" {
		return newDSBolt(path)
	}else{
		return nil, ErrDataEngineNotSupported
	}
}
