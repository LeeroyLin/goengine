package iface

import "go.etcd.io/etcd/clientv3"

type IETCD interface {
	GetClient() *clientv3.Client
	RunDefault(endpoints []string, ttl int64) error
	RunWithConfig(endpoints []string, ttl int64, cfg clientv3.Config) error
	Stop()
	Put(key, value string) error
	Get(key string) (*clientv3.GetResponse, error)
	Delete(key string) error
	Watch(key string, handler func(evt *clientv3.Event))
}

type IETCDGetter interface {
	GetETCD() IETCD
}
