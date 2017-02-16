package server

import (
	"net"
	"bufio"
	// "net/http"
	// "golang.org/x/net/websocket"
	// "github.com/julienschmidt/httprouter"
	"errors"
)

type httpListener struct {
	addr net.Addr
	C chan net.Conn
}

func newHttpListener(addr net.Addr) httpListener {
	return httpListener{
		addr: addr,
		C: make(chan net.Conn),
	}
}

func (l httpListener) Accept() (net.Conn, error) {
	conn, ok := <-l.C
	if !ok {
		return nil, errors.New("listener closed")
	}
	return conn, nil
}

func (l httpListener) Push(rd *bufio.Reader, c net.Conn) {
	l.C <- &httpConn{c, rd}
}

func (l httpListener) Close() error {
	close(l.C)
	return nil
}

func (l httpListener) Addr() net.Addr {
	return l.addr
}


type httpConn struct {
	net.Conn
	rd *bufio.Reader
}

func (hc *httpConn) Read(p []byte) (int, error) {
	return hc.rd.Read(p)
}
