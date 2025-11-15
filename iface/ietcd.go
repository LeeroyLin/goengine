package iface

import "github.com/coreos/etcd/clientv3"

type IETCD interface {
	GetClient() *clientv3.Client
	RunDefault(endpoints []string, ttl int64, action func())
	RunWithConfig(endpoints []string, ttl int64, cfg clientv3.Config, action func())
	Stop()
	Put(key, value string) error
	Get(key string) (*clientv3.GetResponse, error)
	Delete(key string) error
	Watch(key string, handler func(evt *clientv3.Event))
}
