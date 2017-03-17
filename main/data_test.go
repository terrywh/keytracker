package main

import (
	"testing"
)

func TestDataKey(t *testing.T) {
	key := DataKey("abc_")
	if len(key) != 14 {
		t.Fatal("failed DataKey:", key)
	}
}