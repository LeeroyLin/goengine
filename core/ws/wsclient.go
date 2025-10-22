package ws

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface/inetwork"
	"github.com/LeeroyLin/goengine/iface/iwebsocket"
	"github.com/gorilla/websocket"
	"sync"
)

type WSClient struct {
	Url        string
	dataPack   inetwork.IDataPack
	exitChan   chan interface{}
	msgHandler iwebsocket.IWSMsgHandler
	connMgr    iwebsocket.IWSConnManager
	sync.Mutex
	maxMsgBuffChanLen uint32
	workerPoolSize    uint32
}

func NewWSClient(maxMsgBuffChanLen, workerPoolSize, maxWorkerTaskLen uint32, url string, dataPack inetwork.IDataPack) *WSClient {
	c := &WSClient{
		Url:               url,
		dataPack:          dataPack,
		exitChan:          make(chan interface{}),
		msgHandler:        NewWSMsgHandler(workerPoolSize, maxWorkerTaskLen),
		connMgr:           NewWSConnManager(),
		maxMsgBuffChanLen: maxMsgBuffChanLen,
	}

	return c
}

func (c *WSClient) Start() {
	elog.Info("[WSClient] start connect server: ", c.Url)

	c.exitChan = make(chan interface{})

	go func() {
		// 连接
		conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)
		if err != nil {
			elog.Error("[WSClient] connect server err: ", c.Url, err)
			return
		}

		// 开启工作池
		c.msgHandler.StartWorkerPool()

		dealConn := NewWSConnection(c.workerPoolSize, c.maxMsgBuffChanLen, c, conn, 1, c.msgHandler)

		dealConn.Start()

		select {
		case <-c.exitChan:
			c.connMgr.StopAllConn()

			elog.Infof("[WSClient] %s ws client stoped\n", c.Url)
			return
		}
	}()
}

func (c *WSClient) Stop() {
	select {
	case <-c.exitChan:
		return
	default:
		close(c.exitChan)
	}
}

func (s *WSClient) AddRouter(msgId uint32, router iwebsocket.WSRouterHandler) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *WSClient) SetDefaultRouter(router iwebsocket.WSRouterHandler) {
	s.msgHandler.SetDefaultRouter(router)
}

func (s *WSClient) GetConnMgr() iwebsocket.IWSConnManager {
	return s.connMgr
}

func (s *WSClient) GetDataPack() inetwork.IDataPack {
	return s.dataPack
}

func (s *WSClient) RecycleId(connId uint32) {
}
