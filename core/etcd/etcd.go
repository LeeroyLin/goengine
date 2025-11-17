package etcd

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	connCb func()
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

func (e *ETCD) Run(ttl int64, cfg clientv3.Config, connCb func()) error {
	e.endpoints = cfg.Endpoints
	e.ttl = ttl
	e.connCb = connCb

	var err error
	e.client, err = clientv3.New(cfg)

	if err != nil {
		elog.Error("[ECTD] Run etcd failed. err:", err)
		return err
	}

	elog.Info("[ETCD] init success.")

	// 创建租约
	go e.createLease()

	return nil
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

func (e *ETCD) Put(key, value string, timeout time.Duration, opts ...clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	if len(opts) > 0 {
		allOpts := []clientv3.OpOption{clientv3.WithLease(e.leaseId)}
		allOpts = append(allOpts, opts...)
		_, err := e.client.Put(ctx, key, value, allOpts...)
		cancel()

		return err

	} else {
		_, err := e.client.Put(ctx, key, value)
		cancel()
		return err
	}
}

func (e *ETCD) Get(key string, timeout time.Duration, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	resp, err := e.client.Get(ctx, key, opts...)
	cancel()

	return resp, err
}

func (e *ETCD) Delete(key string, timeout time.Duration, opts ...clientv3.OpOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_, err := e.client.Delete(ctx, key, opts...)
	cancel()
	return err
}

func (e *ETCD) Watch(key string, handler func(evt *clientv3.Event)) {
	go func() {
		select {
		case <-e.closeChan:
			return
		default:
			break
		}

		watchChan := e.client.Watch(context.Background(), key)

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

func (e *ETCD) StartTick(duration time.Duration, handler func()) {
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-e.closeChan:
			return
		case <-ticker.C:
			handler()
		}
	}
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

// 创建租约
func (e *ETCD) createLease() {
	ETCDDelay.Reset()

	for {
		select {
		case <-e.closeChan:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			// 创建租约
			leaseResp, err := e.client.Grant(ctx, e.ttl)
			cancel()

			if err != nil {
				elog.Error("[ETCD] create lease failed. retry later. err:", err)
				ETCDDelay.Delay()
				continue
			}

			elog.Info("[ETCD] create lease success.")
			e.leaseId = leaseResp.ID

			e.RLock()
			cb := e.connCb
			e.RUnlock()
			if cb != nil {
				cb()
			}

			// 开始续租
			e.startKeepAlive()

			ETCDDelay.Reset()
		}
	}
}

// 开始续租
func (e *ETCD) startKeepAlive() {
	ticker := time.NewTicker(time.Duration(e.ttl-5) * time.Second)

	for {
		select {
		case <-e.closeChan:
			// 取消租约
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := e.client.Revoke(ctx, e.leaseId)
			cancel()
			if err != nil {
				elog.Fatal("[ETCD] revoke lease failed. err:", err)
			}
			return
		case <-ticker.C:
			// 续租
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, err := e.client.KeepAliveOnce(ctx, e.leaseId)
			cancel()
			if err != nil {
				elog.Error("[ETCD] keep alive failed. err:", err)
				return
			}
		}
	}
}
