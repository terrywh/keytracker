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
	if (r.X & 0x01) != 0 { // 数据设置
		if (r.X & 0x04) != 0 { // 后缀
			r.K = DataKey(r.K)
			DataWrite(s, r.K, r.V, /* y=*/1) // y=1 后缀推送
		}
		if (r.X & 0x02) == 0 {
			s.AddTag(Tag{r.K, false}) // 临时数据需要在 Close 时删除
		}
		if DataSet(r.K, r.V) { // 数据发生变更
			WatcherNotify(r.K, r.V)
		}
	} else if (r.X & 256) != 0 { // 监控
		WatcherAppend(r.K, s)
	} else if (r.X & 512) != 0 {
		DataGet(r.K, s)
	} else if (r.X & 1024) != 0 {
		DataList(r.K, s)
	}
}

func (sh *SessionHandler) CloseHandler(s *server.Session) {
	log.Println("[info] session closed:", s.RemoteAddr)
	atomic.AddInt32(&sessions, -1)
	WatcherCleanup(s)
	DataCleanup(s)
}
