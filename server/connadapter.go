package server

import(
	"github.com/gorilla/websocket"
	"io"
	"errors"
)

var ErrNotTextMessage = errors.New("only message of type 'text' is allowed")

type scTCPSocket struct {
	rd io.Reader
	wt io.Writer
	cs io.Closer
}
func (c *scTCPSocket) Read(p []byte) (int, error) {
	return c.rd.Read(p)
}
func (c *scTCPSocket) Write(p []byte) (int, error) {
	return c.wt.Write(p)
}
func (c *scTCPSocket) Close() error {
	return c.cs.Close()
}

type scWebSocket struct {
	conn *websocket.Conn
}
func (c *scWebSocket) Read(p []byte) (int, error) {
	mt, rd, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	if mt != websocket.TextMessage {
		return 0, ErrNotTextMessage
	}
	return rd.Read(p)
}
func (c *scWebSocket) Write(p []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.TextMessage, p)
	if err == nil {
		return len(p), nil
	}else{
		return 0, err
	}
}
func (c *scWebSocket) Close() error {
	return c.conn.Close()
}
