package app

import "github.com/LeeroyLin/goengine/iface"

func (a *App) GetETCD() iface.IETCD {
	return a.ETCD
}
