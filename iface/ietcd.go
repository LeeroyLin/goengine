package iface

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type IETCD interface {
	GetClient() *clientv3.Client
	Run(ttl int64, cfg clientv3.Config, connCb func()) error
	Stop()
	Put(key, value string, timeout time.Duration, opts ...clientv3.OpOption) error
	Get(key string, timeout time.Duration, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Delete(key string, timeout time.Duration, opts ...clientv3.OpOption) error
	Watch(key string, handler func(evt *clientv3.Event), opts ...clientv3.OpOption)
	WithLease() clientv3.OpOption
}

type IETCDGetter interface {
	GetETCD() IETCD
}
