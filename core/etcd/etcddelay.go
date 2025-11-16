package etcd

import (
	"github.com/LeeroyLin/goengine/core/utils"
	"time"
)

const (
	Start = 500 * time.Millisecond
	Max   = 4 * time.Second
	Multi = 2
)

var ETCDDelay = utils.NewPowDelay(Start, Max, Multi)
