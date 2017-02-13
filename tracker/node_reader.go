package tracker

import (
	"time"
	"container/list"
	"encoding/json"
	"log"
	"fmt"
)

func nodeReader(e *list.Element) {
	n := e.Value.(*_node)
	// 首次超时设置默认为 15秒
	n.to = 15 * time.Second
	n.expire = time.AfterFunc(n.to, func() {
		nodeTimeout(n)
	})
	d := json.NewDecoder(n.conn)
	s := make(map[string]interface{})
	j := false
	for d.More() {
		if nil != d.Decode(&s) {
			break
		}
		var ok bool
		ok = true
		switch action := s["action"]; action {
		case "hello":
			ok = nodeHello(n, s)
			j  = ok
			if j {
				nodeJoined(e)
			}
		case "ping":
			if j {
				ok = nodePing(n, s)
			}
		case "data":
			ok = j
		default:
			ok = false
		}
		if ok && j {
			nodeDataStats(n, s)
			nodeDataDelta(n, s)
			nodeTimeoutReset(n)
		} else {
			log.Printf("[warning] node '%s:%s' from '%s' action '%s' failed.", n.ns, n.key, n.data["remote_addr"], s["action"])
			break
		}
	}
	if j {
		nodeLeaved(e)
	}

	nlock.Lock() // nodes.Remove 操作，需要上锁
	defer nlock.Unlock()
	nodes.Remove(e)
	n.conn.Close() // c.CloseWrite()
	n.expire.Stop()
}

func nodeTimeoutReset(n *_node) {
	// 重用 timer 参考 https://golang.org/pkg/time/#Timer.Reset
	if !n.expire.Stop() {
		<-n.expire.C
	}
	// 使用超时时间设置，可以通过应用端更改
	n.expire.Reset(n.to)
}

func nodeTimeout(n *_node) {
	log.Printf("[warning] node '%s:%s' from '%s' timeout thus being dropped.", n.ns, n.key, n.data["remote_addr"])
	n.conn.CloseRead() // 使 nodeReader 结束（进而清理资源）
}

func nodeHello(n *_node, s map[string]interface{}) bool {
	var ok bool
	n.wlock.Lock()
	defer n.wlock.Unlock()
	if n.ns, ok = s["ns"].(string); !ok {
		n.ns = "default"
	}
	if n.key, ok = s["key"].(string); !ok {
		n.key = newKey()
	}
	nlock.RLock()
	defer nlock.RUnlock()
	for x := nodes.Front(); x != nil; x = x.Next() {
		y := x.Value.(*_node)
		if y != n && y.ns == n.ns && y.key == n.key {
			return false
		}
	}
	if to, ok := s["ping"].(float64); ok && to > 10 { // 最小 ping 间隔 10 秒
		n.to = time.Second * time.Duration(to + 5) // 实际超时检查间隔适当放宽
	}
	log.Printf("[info] node '%s:%s' connected from '%s'", n.ns, n.key, n.data["remote_addr"])
	return true
}

func nodePing(n *_node, s map[string]interface{}) bool {
	n.wlock.Lock()
	defer n.wlock.Unlock()
	n.conn.Write([]byte("{\"action\":\"pong\"}\n"))
	return true
}

func nodeJoined(e *list.Element) {
	n := e.Value.(*_node)
	i := 0

	n.wlock.Lock()
	defer n.wlock.Unlock()
	nlock.RLock()
	defer nlock.RUnlock()
	for x := nodes.Front(); x != nil; x = x.Next() {
		y := x.Value.(*_node)
		if x != e && y.ns == n.ns {
			// 向当前结点报告已连接节点信息
			if i == 0 {
				fmt.Fprintf(n.conn, "{\"action\":\"node\",\"node\":{\"%s\":true", y.key)
			} else {
				fmt.Fprintf(n.conn, ",\"%s\":true", y.key)
			}
			i++
			// 向其他节点报告本节点
			y.wlock.Lock()
			fmt.Fprintf(y.conn, "{\"action\":\"node\",\"node\":{\"%s\":true}}\n", n.key)
			y.wlock.Unlock()
		}
	}
	if i > 0 {
		n.conn.Write([]byte("}}\n"))
	}
}

func nodeLeaved(e *list.Element) {
	nlock.RLock()
	defer nlock.RUnlock()
	n := e.Value.(*_node)
	for x := nodes.Front(); x != nil; x = x.Next() {
		y := x.Value.(*_node)
		if x != e && y.ns == n.ns {
			y.wlock.Lock()
			fmt.Fprintf(y.conn, "{\"action\":\"node\",\"node\":{\"%s\":false}}\n", n.key)
			y.wlock.Unlock()
		}
	}
}

func nodeDataStats(n *_node, s map[string]interface{}) {
	n.wlock.Lock()
	defer n.wlock.Unlock()
	if data, ok := s["stats"].(map[string]interface{}); ok {
		for k, v := range data {
			n.data[k] = v
		}
	}
}

func nodeDataDelta(n *_node, s map[string]interface{}) {
	n.wlock.Lock()
	defer n.wlock.Unlock()
	if data, ok := s["delta"].(map[string]interface{}); ok {
		for k, v := range data {
			v1, _ := n.data[k].(int64)
			v2, _ := v.(float64)
			n.data[k] = int64(v2) + v1
		}
	}
}
