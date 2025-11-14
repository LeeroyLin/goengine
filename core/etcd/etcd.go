package etcd

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

type ETCD struct {
	Endpoints []string
	Client    *clientv3.Client
	LeaseId   clientv3.LeaseID
	closeChan chan interface{}
	ttl       int64
	sync.RWMutex
}

func NewETCD(endpoints []string, ttl int64) *ETCD {
	e := &ETCD{
		Endpoints: endpoints,
		closeChan: make(chan interface{}),
		ttl:       ttl,
	}

	return e
}

func (e *ETCD) RunDefault(action func()) {
	go func() {
		e.Lock()

		var err error
		e.Client, err = clientv3.New(clientv3.Config{
			Endpoints:   e.Endpoints,
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

func (e *ETCD) RunWithConfig(cfg clientv3.Config, action func()) {
	go func() {
		e.Lock()

		var err error
		e.Client, err = clientv3.New(cfg)

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
	_, err := e.Client.Put(context.Background(), key, value, clientv3.WithLease(e.LeaseId))
	return err
}

func (e *ETCD) Get(key string) (*clientv3.GetResponse, error) {
	return e.Client.Get(context.Background(), key)
}

func (e *ETCD) Delete(key string) error {
	_, err := e.Client.Delete(context.Background(), key)
	return err
}

func (e *ETCD) Watch(key string, handler func(evt *clientv3.Event)) {
	go func() {
		e.RLock()
		if e.Client == nil {
			return
		}

		watchChan := e.Client.Watch(context.Background(), key)
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

	if e.Client != nil {
		err := e.Client.Close()
		if err != nil {
			elog.Fatal("[ETCD] close failed. err:", err)
			return
		}
		e.Client = nil
	}
}

func (e *ETCD) createLease() {
	// 创建租约
	leaseResp, err := e.Client.Grant(context.Background(), e.ttl)
	if err != nil {
		elog.Error("[ETCD] create lease failed. err:", err)
		return
	}
	e.LeaseId = leaseResp.ID

	ticker := time.NewTicker(time.Duration(e.ttl-5) * time.Second)

	for {
		select {
		case <-e.closeChan:
			// 取消租约
			_, err := e.Client.Revoke(context.Background(), e.LeaseId)
			if err != nil {
				elog.Fatal("[ETCD] revoke lease failed. err:", err)
			}
			return
		case <-ticker.C:
			// 续租
			_, err := e.Client.KeepAliveOnce(context.Background(), e.LeaseId)
			if err != nil {
				elog.Error("[ETCD] keep alive failed. err:", err)
			}
		}
	}
}
