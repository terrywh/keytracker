package server

import (
	"io"
	"sync"
	"encoding/json"
	"bufio"
	"net"
	"github.com/gorilla/websocket"
	"errors"
	"container/list"
)

type Request struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
	X int         `json:"x"`
}

type sessionConnPeeker struct {
	rd   *bufio.Reader
	conn net.Conn
}

func (c *sessionConnPeeker) Read(p []byte) (int, error) {
	return c.rd.Read(p)
}
func (c *sessionConnPeeker) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

func (c *sessionConnPeeker) Close() error {
	return c.conn.Close()
}

type sessionConnWebsocket struct {
	conn *websocket.Conn
}

func WrapWebsocket(conn *websocket.Conn) io.ReadWriteCloser {
	return sessionConnWebsocket{conn}
}

func (c sessionConnWebsocket) Read(p []byte) (int, error) {
	mt, rd, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	if mt != websocket.TextMessage {
		return 0, errors.New("only TextMessage is allowed")
	}
	return rd.Read(p)
}

func (c sessionConnWebsocket) Write(p []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.TextMessage, p)
	if err == nil {
		return len(p), nil
	}else{
		return 0, err
	}
}

func (c sessionConnWebsocket) Close() error {
	return c.conn.Close()
}


type SessionHandler interface {
	StartHandler(s *Session)
	RequestHandler(s *Session, r *Request)
	CloseHandler(s *Session)
}

type Session struct {
	conn io.ReadWriteCloser
	lock *sync.RWMutex
	RemoteAddr string
	tags *list.List
}

func NewSession(conn io.ReadWriteCloser, addr string) *Session {
	return &Session{
		conn: conn,
		lock: &sync.RWMutex{},
		RemoteAddr: addr,
		tags: list.New(),
	}
}

func (s *Session) Start(sh SessionHandler) {
	d := json.NewDecoder(s.conn)
	for d.More() {
		r := Request{}
		if err := d.Decode(&r); err != nil || r.K == "" {
			continue
		}
		sh.RequestHandler(s, &r)
	}
	s.conn.Close()
	sh.CloseHandler(s)
}
// func (s *Session) Close() {
// 	s.conn.Close()
// }
func (s *Session) Write(b []byte) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.conn.Write(b)
}

func (s *Session) AddTag(tag interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.tags.PushBack(tag)
}

func (s *Session) WalkTag(cb func(interface{}) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for tag := s.tags.Front(); tag != nil; tag = tag.Next() {
		if !cb(tag.Value) {
			break
		}
	}
}
