package main

import (
	"github.com/terrywh/keytracker/server"
	"sync/atomic"
	"log"
)

type Tag struct {
	Key string
	IsWatcher bool
}

type SessionHandler struct {
}

var handler  *SessionHandler

func init() {
	handler  = &SessionHandler{}
}

var sessions int32

func (sh *SessionHandler) StartHandler(s *server.Session) {
	log.Println("[info] session started:", s.RemoteAddr)
	atomic.AddInt32(&sessions, 1)
}

func (sh *SessionHandler) RequestHandler(s *server.Session, r *server.Request) {
	r.K = DataKeyFlat(r.K)
	if (r.X & 0x01) != 0 { // 数据设置
		if (r.X & 0x04) != 0 { // 后缀
			r.K = DataKey(r.K)
			DataWrite(s, r.K, r.V, /* y=*/1) // y=1 后缀推送
		}
		if (r.X & 0x02) == 0 {
			s.AddElement(r.K) // 临时数据需要在 Close 时删除
		}
		if DataSet(r.K, r.V) { // 数据发生变更
			WatcherNotify(r.K, r.V, r.K)
		}
	} else if (r.X & 256) != 0 { // 监控
		var v, _ = r.V.(float64)

		if int(v) == 0 {
			WatcherRemove(r.K, s)
		} else if int(v) == 1 {
			DataList(r.K, s, 2, func() {
				WatcherAppend(r.K, s) // 防止 watcher 和 list 之间空隙时数据的变更
				s.AddWatcher(r.K)
			})
		}
	} else if (r.X & 512) != 0 {
		DataGet(r.K, s)
	} else if (r.X & 1024) != 0 {
		DataList(r.K, s, 0, nil)
	} else { // 删除数据
		DataDel(r.K)
	}
}

func (sh *SessionHandler) CloseHandler(s *server.Session) {
	log.Println("[info] session closed:", s.RemoteAddr)
	atomic.AddInt32(&sessions, -1)
	WatcherCleanup(s)
	DataCleanup(s)
}
