package trie

import (
	"testing"
)



func TestTrieGetSet(t *testing.T) {
	n := NewTrie()
	key1 := "/a/b/c"
	val1 := "1111111111111111"
	key2 := "/a/b"
	val2 := "22222222222"
	key3 := "/a"
	val3 := "3333333"
	n.Create(key1).SetValue(val1)
	n.Create(key2).SetValue(val2)
	n.Create(key3).SetValue(val3)

	if n.Get(key1).GetValue() != val1 || n.Get(key2).GetValue() != val2 || n.Get(key3).GetValue() != val3 {
		t.Fail()
	}

	if n.Get(key1).Key != "c" || n.Get(key2).Key != "b" || n.Get(key3).Key != "a" {
		t.Log(n.Get(key1).Key, key1)
		t.Fail()
	}
}

func TestTrieRemove(t *testing.T) {
	n := NewTrie()
	key1 := "/a/b/c"
	val1 := "1111111111111111"
	key2 := "/a/b"
	val2 := "22222222222"
	key3 := "/a"
	val3 := "3333333"
	n.Create(key1).SetValue(val1)
	n.Create(key2).SetValue(val2)
	n.Create(key3).SetValue(val3)

	n.Remove(key1)
	if n.Get(key1) != nil {
		t.Log(n.Get(key1))
		t.Fail()
	}
	n.Remove(key2)
	if n.Get(key2) != nil {
		t.Log(n.Get(key2))
		t.Fail()
	}
	n.Remove(key3)
	if n.Get(key3) != nil {
		t.Log(n.Get(key3))
		t.Fail()
	}
}

// var n *Node = NewNode()
//
// func BenchmarkTrieSet(b *testing.B) {
// 	b.N = 9000000
// 	x := 1
// 	y := 1
// 	z := 1
// 	o := 1
// 	p := 1
// 	q := 1
// 	for i:=0; i<b.N; i++ {
// 		k := []string{
// 			string(x%100000),
// 			string(y%100000),
// 			string(z%100000),
// 			// string(o%100000),
// 			// string(p%100000),
// 			// string(q%100000),
// 		}
// 		n.Create(k)
// 		x += 1
// 		y += 2
// 		z += 3
// 		o += 4
// 		p += 5
// 		q += 6
// 	}
// }
//
// func BenchmarkTrieGet(b *testing.B) {
// 	b.N = 9000000
// 	x := 2
// 	y := 2
// 	z := 2
// 	o := 2
// 	p := 2
// 	q := 2
// 	for i:=0; i<b.N; i++ {
// 		k := []string{
// 			string(x%100000),
// 			string(y%100000),
// 			string(z%100000),
// 			// string(o%100000),
// 			// string(p%100000),
// 			// string(q%100000),
// 		}
// 		_ = n.Get(k)
// 		x += 1
// 		y += 2
// 		z += 3
// 		o += 4
// 		p += 5
// 		q += 6
// 	}
// }

//
// var m map[string]interface{} = make(map[string]interface{})
//
// func BenchmarkMapSet(b *testing.B) {
// 	x := 1
// 	y := 1
// 	z := 1
// 	o := 1
// 	p := 1
// 	q := 1
// 	for i:=0; i<b.N; i++ {
// 		k := ""
// 		k += string(x)
// 		k += string(y)
// 		k += string(z)
// 		// k += string(o)
// 		// k += string(p)
// 		// k += string(q)
//
// 		m[k] = "true"
// 		x++
// 		y++
// 		z++
// 		o++
// 		p++
// 		q++
// 	}
// }
//
// func BenchmarkMapGet(b *testing.B) {
// 	x := 1
// 	y := 1
// 	z := 1
// 	o := 1
// 	p := 1
// 	q := 1
// 	for i:=0; i<b.N; i++ {
// 		k := ""
// 		k += string(x)
// 		k += string(y)
// 		k += string(z)
// 		// k += string(o)
// 		// k += string(p)
// 		// k += string(q)
//
// 		_ = m[k]
// 		x++
// 		y++
// 		z++
// 		o++
// 		p++
// 		q++
// 	}
// }
//
// func BenchmarkMapIteration(b *testing.B) {
// 	for i:=0; i<b.N; i++ {
// 		for _,_ = range m {
//
// 		}
// 	}
// }
