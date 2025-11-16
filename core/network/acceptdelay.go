package network

import (
	"github.com/LeeroyLin/goengine/core/utils"
	"time"
)

const (
	Start = 5 * time.Millisecond
	Max   = 1 * time.Second
	Multi = 2
)

var AcceptDelay = utils.NewPowDelay(Start, Max, Multi)
