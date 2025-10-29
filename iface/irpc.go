package iface

import (
	"google.golang.org/grpc"
)

type IRPC interface {
	NewServer(url string, opt ...grpc.ServerOption) *grpc.Server
	RemoveServer(url string)
	StartServe()

	NewClientConn(url string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	RemoveClientConn(url string)

	ClearAll()
}

type IRPCGetter interface {
	GetRPC() IRPC
}
