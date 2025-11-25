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
	_, err := e.client.Put(ctx, key, value, opts...)
	cancel()

	return err
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

func (e *ETCD) Watch(key string, handler func(evt *clientv3.Event), opts ...clientv3.OpOption) {
	go func() {
		delay := NewETCDWatchDelay()

		for {
			select {
			case <-e.closeChan:
				return
			default:
				watchChan := e.client.Watch(context.Background(), key, opts...)

				for {
					reWatch := false

					select {
					case <-e.closeChan:
						return
					case resp, ok := <-watchChan:
						if !ok {
							elog.Error("[ETCD] channel closed.")
							reWatch = true
							break
						}
						if resp.Err() != nil {
							elog.Error("[ETCD] watch err.", resp.Err())
							reWatch = true
							break
						}
						for _, event := range resp.Events {
							handler(event)
						}
						delay.Reset()
					}

					if reWatch {
						break
					}
				}

				delay.Delay()
			}
		}
	}()
}

func (a *ETCD) WithLease() clientv3.OpOption {
	return clientv3.WithLease(a.leaseId)
}

func (e *ETCD) doClose() {
	e.Lock()
	defer e.Unlock()

	if e.client != nil {
		// 取消租约
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_, err := e.client.Revoke(ctx, e.leaseId)
		cancel()
		if err != nil {
			elog.Fatal("[ETCD] revoke lease failed. err:", err)
		}

		// 关闭客户端
		err = e.client.Close()
		if err != nil {
			elog.Fatal("[ETCD] close failed. err:", err)
			return
		}
		e.client = nil
	}
}

// 创建租约
func (e *ETCD) createLease() {
	ETCDLeaseDelay.Reset()

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
				ETCDLeaseDelay.Delay()
				continue
			}

			elog.Info("[ETCD] create lease success.")
			e.leaseId = leaseResp.ID

			if e.connCb != nil {
				e.connCb()
			}

			// 开始续租
			e.startKeepAlive()

			ETCDLeaseDelay.Reset()
		}
	}
}

// 开始续租
func (e *ETCD) startKeepAlive() {
	ticker := time.NewTicker(time.Duration(e.ttl-5) * time.Second)

	for {
		select {
		case <-e.closeChan:
			return
		case <-ticker.C:
			// 续租
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			_, err := e.client.KeepAliveOnce(ctx, e.leaseId)
			cancel()
			if err != nil {
				elog.Error("[ETCD] keep alive failed. err:", err)
				return
			}
		}
	}
}
