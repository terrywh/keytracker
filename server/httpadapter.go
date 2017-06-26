package server

import (
	"net"
	"bufio"
	"errors"
)

type httpListener struct {
	addr net.Addr
	C chan net.Conn
	closed bool
}
// 实现
func newHttpListener(addr net.Addr) *httpListener {
	return &httpListener{
		addr: addr,
		C: make(chan net.Conn),
		closed: false,
	}
}
func (l *httpListener) Push(rd *bufio.Reader, c net.Conn) {
	l.C <- &httpConn{c, rd}
}
// 实现 net.Listener 接口
func (l *httpListener) Accept() (net.Conn, error) {
	conn, ok := <-l.C
	if !ok {
		return nil, errors.New("listener closed")
	}
	return conn, nil
}
func (l *httpListener) Close() error {
	if l.closed {
		return nil
	}
	l.closed = true
	close(l.C)
	return nil
}
func (l *httpListener) Addr() net.Addr {
	return l.addr
}

type httpConn struct {
	net.Conn
	rd *bufio.Reader
}
func (hc *httpConn) Read(p []byte) (int, error) {
	return hc.rd.Read(p)
}
