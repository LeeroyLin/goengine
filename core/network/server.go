package network

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/config"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/pool"
	"github.com/LeeroyLin/goengine/iface/inetwork"
	"net"
)

type Server struct {
	IPVersion  string
	IP         string
	Port       int
	conf       *config.ConfNetServicePattern
	msgHandler inetwork.IMsgHandler
	connMgr    inetwork.IConnManager
	dataPack   inetwork.IDataPack
	exitChan   chan interface{}
	idPool     *pool.IdPool[uint32]
}

func (s *Server) Start() {
	s.exitChan = make(chan interface{})

	go func() {
		// 获得addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			elog.Panic("[Server] resolve tcp addr err: ", err)
		}

		// 监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			elog.Panic("[Server] listen TCP err: ", err)
		}

		elog.Infof("[Server] %s server start at port %d\n", s.IP, s.Port)

		// 开启工作池
		s.msgHandler.StartWorkerPool()

		cid := s.idPool.Get()

		go func() {
			for {
				// 如果连接数达上限
				if s.GetConnMgr().Size() >= s.conf.MaxConn {
					// 延时
					AcceptDelay.Delay()
					continue
				}

				// 等待连接
				conn, err := listener.AcceptTCP()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}

					elog.Errorf("[Server] accept tcp err: %v\n", err)
					// 延时
					AcceptDelay.Delay()
					continue
				}

				select {
				case <-s.exitChan:
					// 直接关闭
					err := conn.Close()
					if err != nil {
						elog.Info("[Server] close conn directly err: ", err)
					}
					break
				default:
					AcceptDelay.Reset()

					dealConn := NewConnection(s.conf.WorkerPoolSize, s.conf.MaxMsgBuffChanLen, s, conn, cid, s.msgHandler)
					cid++

					dealConn.Start()
				}
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			elog.Info("[Server] close listener...")
			if err != nil {
				elog.Error("[Server] close listener err: ", err)
			}

			s.connMgr.StopAllConn()

			elog.Infof("[Server] %s server stoped at port %d\n", s.IP, s.Port)
		}
	}()
}

func (s *Server) Stop() {
	select {
	case <-s.exitChan:
		return
	default:
		close(s.exitChan)
	}
}

func (s *Server) AddRouter(msgId uint32, router inetwork.RouterHandler) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) SetDefaultRouter(router inetwork.RouterHandler) {
	s.msgHandler.SetDefaultRouter(router)
}

func (s *Server) GetConnMgr() inetwork.IConnManager {
	return s.connMgr
}

func (s *Server) GetDataPack() inetwork.IDataPack {
	return s.dataPack
}

func (s *Server) RecycleId(connId uint32) {
	s.idPool.Set(connId)
}

func NewServer(conf *config.ConfNetServicePattern, dataPack inetwork.IDataPack) inetwork.IServer {
	s := &Server{
		IPVersion:  conf.IPVersion,
		IP:         conf.IP,
		Port:       conf.Port,
		conf:       conf,
		msgHandler: NewMsgHandler(conf.WorkerPoolSize, conf.MaxWorkerTaskLen),
		connMgr:    NewConnManager(),
		dataPack:   dataPack,
		exitChan:   make(chan interface{}),
		idPool:     pool.NewUint32IdPool(uint32(conf.MaxConn)),
	}

	return s
}
