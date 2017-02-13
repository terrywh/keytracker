package tracker

import (
	"net"
	"container/list"
	"time"
	"sync"
	"github.com/terrywh/ntracker/config"
)

type _node struct {
	conn   *net.TCPConn
	expire *time.Timer
	to      time.Duration
	ns      string
	key     string
	data    map[string]interface{}
	wlock  *sync.Mutex // 写锁定
}
// // 实现 io.Writer 接口
// func (n *_node) Write(p []byte) (int, error) {
// 	// 由于写入动作可能由 api 接口 和 node 接口触发，内部加锁
// 	n.wlock.Lock()
// 	defer n.wlock.Unlock()
//
// 	return n.conn.Write(p)
// }

var (
	server func(l *net.TCPListener)
	nodes *list.List
	nlock *sync.RWMutex
)

func init() {
	server = nodeServerHandler
	nodes  = list.New()
	nlock  = &sync.RWMutex{}
}

func nodeServerHandler(l *net.TCPListener) {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			continue
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(config.NodeKeepAlive)
		n := &_node{
			conn: conn,
			data: make(map[string]interface{}),
			wlock: &sync.Mutex{},
		}
		n.data["remote_addr"] = conn.RemoteAddr().String()

		nlock.Lock()  // nodes 的操作需要加锁
		e := nodes.PushBack(n)
		nlock.Unlock()

		go nodeReader(e)
	}
}

func GetNodeServer() func(l *net.TCPListener) {
	return server
}
