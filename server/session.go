package server

import (
	"io"
	"sync"
	"encoding/json"
)
type Session struct {
	conn io.ReadWriteCloser
	lock *sync.Mutex
	RemoteAddr string
}
func NewSession(conn io.ReadWriteCloser, addr string) *Session {
	return &Session{
		conn: conn,
		lock: &sync.Mutex{},
		RemoteAddr: addr,
	}
}
func (s *Session) Start(svr *Server) {
	svr.OnStart(s)
	d := json.NewDecoder(s.conn)
	for d.More() && !svr.closing {
		r := Request{}
		if err := d.Decode(&r); err != nil || r.K == "" {
			continue
		}
		if !svr.closing {
			svr.OnRequest(s, &r)
		}
	}
	s.Close()
	svr.OnClose(s)
}
func (s *Session) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.conn.Close()
}
func (s *Session) Write(b []byte) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.conn.Write(b)
}
