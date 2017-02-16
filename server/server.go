package server

import (
	"net"
	"log"
	"time"
	"bufio"
	"net/http"
)




func ListenAndServe(addr string, sh SessionHandler, hh http.Handler) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic("[server] failed to listen on " + addr)
	}else{
		log.Println("[info] server started on", addr)
	}

	httpServer := http.Server{Handler: hh}
	httpListener := newHttpListener(l.Addr())

	go httpServer.Serve(httpListener)
	for {
		ccc, err := l.Accept()
		if err != nil {
			log.Println("[warning] failed to accept socket: ", err)
			continue
		}
		conn := ccc.(*net.TCPConn)
		initConn(conn)
		// 至少发送一个字节数据后开始服务（用于识别协议）
		rd := bufio.NewReader(ccc)
		pbyte, _ := rd.Peek(1)
		if pbyte[0] == byte('{') {
			// detect 协议
			s := NewSession(&sessionConnPeeker{rd, conn}, conn.RemoteAddr().String())
			go s.Start(sh)
		}else{
			httpListener.Push(rd, ccc)
		}
	}
}

func initConn(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(15 * time.Second)
	// conn.SetLinger(-1)
}
