package iface

import clientv3 "go.etcd.io/etcd/client/v3"

type IETCD interface {
	GetClient() *clientv3.Client
	Run(ttl int64, cfg clientv3.Config, connCb func()) error
	Stop()
	Put(key, value string) error
	Get(key string) (*clientv3.GetResponse, error)
	Delete(key string) error
	Watch(key string, handler func(evt *clientv3.Event))
}

type IETCDGetter interface {
	GetETCD() IETCD
}
