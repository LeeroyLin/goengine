package etcd

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

type ETCD struct {
	endpoints []string
	client    *clientv3.Client
	leaseId   clientv3.LeaseID
	closeChan chan interface{}
	ttl       int64
	sync.RWMutex
}

func NewETCD() iface.IETCD {
	e := &ETCD{
		closeChan: make(chan interface{}),
	}

	return e
}

func (e *ETCD) GetClient() *clientv3.Client {
	return e.client
}

func (e *ETCD) RunDefault(endpoints []string, ttl int64, action func()) {
	e.endpoints = endpoints
	e.ttl = ttl

	go func() {
		e.Lock()

		var err error
		e.client, err = clientv3.New(clientv3.Config{
			Endpoints:   e.endpoints,
			DialTimeout: 5 * time.Second,
		})

		e.Unlock()

		if err != nil {
			elog.Error("[ECTD] Run etcd failed. err:", err)
			return
		}

		elog.Info("[ETCD] run default success.")

		action()

		// 创建租约
		go e.createLease()

		for {
			select {
			case <-e.closeChan:
				return
			}
		}
	}()
}

func (e *ETCD) RunWithConfig(endpoints []string, ttl int64, cfg clientv3.Config, action func()) {
	e.endpoints = endpoints
	e.ttl = ttl

	go func() {
		e.Lock()

		var err error
		e.client, err = clientv3.New(cfg)

		e.Unlock()

		if err != nil {
			elog.Error("[ECTD] Run etcd failed. err:", err)
			return
		}

		elog.Info("[ETCD] run with config success.")

		action()

		// 创建租约
		go e.createLease()

		for {
			select {
			case <-e.closeChan:
				return
			}
		}
	}()
}

func (e *ETCD) Stop() {
	select {
	case <-e.closeChan:
		return
	default:
		close(e.closeChan)
		e.doClose()
	}
}

func (e *ETCD) Put(key, value string) error {
	_, err := e.client.Put(context.Background(), key, value, clientv3.WithLease(e.leaseId))
	return err
}

func (e *ETCD) Get(key string) (*clientv3.GetResponse, error) {
	return e.client.Get(context.Background(), key)
}

func (e *ETCD) Delete(key string) error {
	_, err := e.client.Delete(context.Background(), key)
	return err
}

func (e *ETCD) Watch(key string, handler func(evt *clientv3.Event)) {
	go func() {
		e.RLock()
		if e.client == nil {
			return
		}

		watchChan := e.client.Watch(context.Background(), key)
		e.RUnlock()

		for {
			select {
			case <-e.closeChan:
				return
			case resp := <-watchChan:
				for _, event := range resp.Events {
					handler(event)
				}
			}
		}
	}()
}

func (e *ETCD) doClose() {
	e.Lock()
	defer e.Unlock()

	if e.client != nil {
		err := e.client.Close()
		if err != nil {
			elog.Fatal("[ETCD] close failed. err:", err)
			return
		}
		e.client = nil
	}
}

func (e *ETCD) createLease() {
	// 创建租约
	leaseResp, err := e.client.Grant(context.Background(), e.ttl)
	if err != nil {
		elog.Error("[ETCD] create lease failed. err:", err)
		return
	}
	e.leaseId = leaseResp.ID

	ticker := time.NewTicker(time.Duration(e.ttl-5) * time.Second)

	for {
		select {
		case <-e.closeChan:
			// 取消租约
			_, err := e.client.Revoke(context.Background(), e.leaseId)
			if err != nil {
				elog.Fatal("[ETCD] revoke lease failed. err:", err)
			}
			return
		case <-ticker.C:
			// 续租
			_, err := e.client.KeepAliveOnce(context.Background(), e.leaseId)
			if err != nil {
				elog.Error("[ETCD] keep alive failed. err:", err)
			}
		}
	}
}
