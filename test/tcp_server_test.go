package test

import (
	"github.com/LeeroyLin/goengine/core/config"
	"github.com/LeeroyLin/goengine/core/network"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	c := &Conf{
		ConfBase:              config.NewConfBase(),
		ConfNetServicePattern: config.NewConfNetServicePattern(),
	}
	c.Setup(c, "")
	s := network.NewServer(&c.ConfNetServicePattern, network.NewDataPack(c.MaxPacketSize))

	s.Start()

	//time.Sleep(time.Second * 2)
	//s.Stop()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
