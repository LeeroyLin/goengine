package test

import (
	"engine/core/conf"
	"engine/core/elog"
	"engine/core/network"
	"engine/core/ws"
	"engine/iface/iwebsocket"
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

func TestWSServer(t *testing.T) {
	c := conf.NewConf()
	c.LoadFromFile("")
	s := ws.NewWSServer(c, network.NewDataPack(c.MaxPacketSize))

	s.AddRouter(1, onMsg1)

	s.Start()

	//time.Sleep(time.Second * 2)
	//s.Stop()

	select {
	case <-time.After(1 * time.Hour):
		break
	}
}
