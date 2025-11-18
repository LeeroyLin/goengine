package etcd

import (
	"github.com/LeeroyLin/goengine/core/utils"
	"time"
)

const (
	LeaseStart = 500 * time.Millisecond
	LeaseMax   = 4 * time.Second
	LeaseMulti = 2
	WatchStart = 1 * time.Second
	WatchMax   = 4 * time.Second
	WatchMulti = 2
)

var ETCDLeaseDelay = utils.NewPowDelay(LeaseStart, LeaseMax, LeaseMulti)

func NewETCDWatchDelay() *utils.PowDelay {
	return utils.NewPowDelay(WatchStart, WatchMax, WatchMulti)
}
