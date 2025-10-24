package test

import (
	"github.com/LeeroyLin/goengine/core/conf"
	"github.com/LeeroyLin/goengine/core/network"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	c := config.NewConf()
	c.LoadFromFile("")
	s := network.NewServer(c, network.NewDataPack(c.MaxPacketSize))

	s.Start()

	//time.Sleep(time.Second * 2)
	//s.Stop()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
