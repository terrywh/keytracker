package trie

import (
	// "sync"
	"strings"
)

type Node struct {
	Key   string // Key 不会修改
	val   interface{}
	child map[string]*Node
	count uint16
// 	lock *sync.RWMutex
}

type Trie struct {
	root *Node
}

func NewTrie() Trie {
	return Trie{
		root: newNode(""),
	}
}

func (t Trie) Create(key string) *Node {
	v1 := t.root
	v2 := t.root
	var ok bool
	for _, k := range strings.Split(key, "/") {
		v1 = v2
		if k == "" {
			continue
		}
		// v1.lock.Lock()
		v2, ok = v1.child[k]
		if !ok {
		 	v2 = newNode(k)
		 	v1.child[k] = v2
			v1.count++
		}
		// v1.lock.Unlock()
	}
	return v2
}

func (t Trie) Remove(key string) *Node {
	v1 := t.root
	v2 := t.root
	keys := strings.Split(key, "/")
	var ok bool
	for _, k := range keys {
		v1 = v2
		if k == "" {
			continue
		}
		// v1.lock.RLock()
		v2, ok = v1.child[k]
		// v1.lock.RUnlock()
		if !ok {
		 	return nil
		}
	}
	// v1.lock.Lock()
	// defer v1.lock.Unlock()
	delete(v1.child, keys[len(keys)-1])
	return v2
}

// 返回数据若为指针类型，可以进行修改，否则请使用 Set 方法通过上层节点设置值
func (t Trie) Get(key string) *Node {
	v1 := t.root
	v2 := t.root
	var ok bool
	for _, k := range strings.Split(key, "/") {
		v1 = v2
		if k == "" {
			continue
		}
		// v1.lock.RLock()
		v2, ok = v1.child[k]
		// v1.lock.RUnlock()
		if !ok {
			return nil
		}

	}
	return v2
}

func newNode(key string) *Node {
	return &Node{
		Key: key,
		val: nil,
		child: map[string]*Node{},
		count: 0,
		// lock: &sync.RWMutex{},
	}
}

func (n *Node) GetValue() interface{} {
	// n.lock.RLock()
	// defer n.lock.RUnlock()
	return n.val
}

func (n *Node) SetValue(val interface{}) bool {
	// n.lock.Lock()
	// defer n.lock.Unlock()
	if n.val != val {
		n.val = val
		return true
	}
	return false
}

func (n *Node) Walk(cb func(*Node) bool) {
	// n.lock.RLock()
	// defer n.lock.RUnlock()
	for _, v := range n.child {
		if !cb(v) {
			break
		}
	}
}
