package tracker

import (
	"crypto/rand"
	"encoding/hex"
)

func newKey() string {
	key := make([]byte, 8)
	n,err:= rand.Read(key)
	if n != 8 || err != nil {
		panic("failed to create node key")
	}
	return hex.EncodeToString(key)
}
