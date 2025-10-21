package test

import (
	"engine/core/conf"
	"engine/core/network"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	c := conf.NewConf()
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
