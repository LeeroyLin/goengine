package test

import (
	"github.com/LeeroyLin/goengine/core/config"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/network"
	"github.com/LeeroyLin/goengine/core/ws"
	"github.com/LeeroyLin/goengine/iface/iwebsocket"
	"testing"
	"time"
)

func onMsg1(req iwebsocket.IWSRequest) {
	elog.Info("Rec ", string(req.GetData()))

	str := "back " + string(req.GetData())
	err := req.GetConnection().SendBuffMsg(2, []byte(str))
	if err != nil {
		elog.Error("send back str err.", err)
		return
	}
}

type Conf struct {
	config.ConfBase
	config.ConfNetServicePattern
}

func TestWSServer(t *testing.T) {
	c := &Conf{
		ConfBase:              config.NewConfBase(),
		ConfNetServicePattern: config.NewConfNetServicePattern(),
	}
	c.Setup(c, "")
	s := ws.NewWSServer(&c.ConfNetServicePattern, network.NewDataPack(c.MaxPacketSize))

	s.AddRouter(1, onMsg1)

	s.Start()

	//time.Sleep(time.Second * 2)
	//s.Stop()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
