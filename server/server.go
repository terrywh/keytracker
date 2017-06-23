package server

import (
	"net"
	"time"
	"bufio"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/websocket"
)

type Request struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
	X int         `json:"x"`
}

type Server struct {
	Router *httprouter.Router
	tcpLst  net.Listener
	htpSvr  http.Server
	htpLst  httpListener
	wskUpg  websocket.Upgrader
	OnStart   func(s *Session)
	OnRequest func(s *Session, r *Request)
	OnClose   func(s *Session)
}

func New() *Server {
	router := httprouter.New()
	svr := Server {
		Router: router,
		htpSvr: http.Server { Handler: router },
		wskUpg: websocket.Upgrader {
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
	router.GET("/session", svr.shWebsocket)
	return &svr;
}

func (svr *Server) ListenAndServe(addr string) {
	var err error
	// TCP 监听
	svr.tcpLst, err = net.Listen("tcp", addr)
	if err != nil {
		// log.Fatal("[fatal] server failed to listen on", addr)
	}
	// log.Println("[info] server started on", addr)
	// HTTP 服务
	svr.htpLst = newHttpListener(svr.tcpLst.Addr())
	go svr.htpSvr.Serve(svr.htpLst)

	// 连接识别逻辑
	for {
		ccc, err := svr.tcpLst.Accept()
		if err != nil {
			// log.Println("[error] failed to accept socket: ", err)
			break
		}
		// 连接基础处理
		conn := ccc.(*net.TCPConn)
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(15 * time.Second)
		cbuf := bufio.NewReader(conn)
		// 协议识别
		cbeg, err := cbuf.Peek(1)
		if err != nil {
			// 发生错误即无法识别协议
			// log.Println("[warning] socket diconnected before protocol detection")
		} else if cbeg[0] == byte('{') {
			// TCP 连接 JSON 数据行协议
			s := NewSession(&scTCPSocket{cbuf, conn, conn}, conn.RemoteAddr().String())
			go s.Start(svr)
		}else{
			// HTTP 协议（可 Upgrade 为 WebSocket）
			svr.htpLst.Push(cbuf, conn)
		}
	}
}

func (svr *Server) shWebsocket(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := svr.wskUpg.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	ss := NewSession(&scWebSocket{conn}, conn.RemoteAddr().String())
	// HTTP 本身已经位于独立的协程中，不需要再次启动新的协程
	ss.Start(svr)
}


func (svr *Server) Close() {
	svr.tcpLst.Close()
	svr.htpSvr.Close()
}
