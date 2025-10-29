package app

import "github.com/LeeroyLin/goengine/iface"

func (a *App) GetRPC() iface.IRPC {
	return a.RPC
}
