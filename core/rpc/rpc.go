package rpc

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/syncmap"
	"google.golang.org/grpc"
	"net"
)

type RPC struct {
	servers     *syncmap.SyncMap[string, *grpc.Server]
	clientConns *syncmap.SyncMap[string, *grpc.ClientConn]
}

func NewRPC() *RPC {
	rpc := &RPC{
		servers:     syncmap.NewSyncMap[string, *grpc.Server](),
		clientConns: syncmap.NewSyncMap[string, *grpc.ClientConn](),
	}

	return rpc
}

func (rpc *RPC) NewServer(url string, opt ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opt...)

	// 添加
	rpc.servers.Add(url, s)

	return s
}

func (rpc *RPC) RemoveServer(url string) {
	s, ok := rpc.servers.Get(url)

	if ok {
		s.Stop()
		rpc.servers.Delete(url)
	}
}

func (rpc *RPC) StartServe() {
	rpc.servers.Range(func(url string, s *grpc.Server) bool {
		rpc.serveOne(url, s)
		return true
	})
}

func (rpc *RPC) serveOne(url string, s *grpc.Server) {
	go func() {
		defer rpc.servers.Delete(url)

		listen, err := net.Listen("tcp", url)
		if err != nil {
			elog.Error("[RPC] listen tcp url err.", url, err)
			return
		}

		elog.Error("[RPC] start serve.", url)

		err = s.Serve(listen)

		if err != nil {
			elog.Error("[RPC] serve err.", url, err)
			return
		}
	}()
}

func (rpc *RPC) NewClientConn(url string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(url, opts...)

	if err != nil {
		return nil, err
	}

	rpc.clientConns.Add(url, conn)

	return conn, nil
}

func (rpc *RPC) RemoveClientConn(url string) {
	c, ok := rpc.clientConns.Get(url)

	if ok {
		c.Close()
		rpc.clientConns.Delete(url)
	}
}

func (rpc *RPC) ClearAll() {
	rpc.servers.Range(func(url string, s *grpc.Server) bool {
		s.Stop()
		return true
	})

	rpc.servers.Clear()

	rpc.clientConns.Range(func(url string, conn *grpc.ClientConn) bool {
		conn.Close()
		return true
	})

	rpc.clientConns.Clear()
}
