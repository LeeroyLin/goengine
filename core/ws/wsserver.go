package ws

import (
	"engine/core/conf"
	"engine/core/elog"
	"engine/iface/inetwork"
	"engine/iface/iwebsocket"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type WSServer struct {
	IPVersion  string
	IP         string
	Port       int
	Url        string
	conf       *conf.Conf
	msgHandler iwebsocket.IWSMsgHandler
	connMgr    iwebsocket.IWSConnManager
	dataPack   inetwork.IDataPack
	exitChan   chan interface{}
	upgrader   websocket.Upgrader
	cid        uint32
}

func (s *WSServer) Start() {
	logStr := s.conf.GetLogStr()
	elog.Infof("[Server] start server. conf: %s", logStr)

	s.exitChan = make(chan interface{})

	// 升级器，用于将 HTTP 连接升级为 WebSocket 连接
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	go func() {
		// 设置 WebSocket 处理路径
		http.HandleFunc("/ws", s.serveWS)

		elog.Infof("[Server] %s server start...\n", s.Url)

		// 开启工作池
		s.msgHandler.StartWorkerPool()

		// 启动服务器
		err := http.ListenAndServe(s.Url, nil)
		if err != nil {
			elog.Panic("[Server] start server err: ", err)
		}

		select {
		case <-s.exitChan:
			s.connMgr.StopAllConn()

			elog.Infof("[Server] %s server stoped at port %d\n", s.IP, s.Port)
		}
	}()
}

func (s *WSServer) Stop() {
	select {
	case <-s.exitChan:
		return
	default:
		close(s.exitChan)
	}
}

func (s *WSServer) AddRouter(msgId uint32, router iwebsocket.WSRouterHandler) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *WSServer) SetDefaultRouter(router iwebsocket.WSRouterHandler) {
	s.msgHandler.SetDefaultRouter(router)
}

func (s *WSServer) GetConnMgr() iwebsocket.IWSConnManager {
	return s.connMgr
}

func (s *WSServer) GetDataPack() inetwork.IDataPack {
	return s.dataPack
}

func NewWSServer(conf *conf.Conf, dataPack inetwork.IDataPack) *WSServer {
	s := &WSServer{
		IPVersion:  conf.IPVersion,
		IP:         conf.IP,
		Port:       conf.Port,
		Url:        fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		conf:       conf,
		msgHandler: NewWSMsgHandler(conf.WorkerPoolSize, conf.MaxWorkerTaskLen),
		connMgr:    NewWSConnManager(),
		dataPack:   dataPack,
		exitChan:   make(chan interface{}),
		cid:        0,
	}

	return s
}

func (s *WSServer) serveWS(w http.ResponseWriter, r *http.Request) {
	// 是否连接数达上限
	if s.connMgr.Len() >= s.conf.MaxConn {
		elog.Error("[Server] Already max conn.")
		return
	}

	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		elog.Error("[Server] Server websocket upgrade err. ", err)
		return
	}

	dealConn := NewWSConnection(s.conf.WorkerPoolSize, s.conf.MaxMsgBuffChanLen, s, conn, s.cid, s.msgHandler)

	// todo 临时用递增id测试
	s.cid++

	dealConn.Start()
}
