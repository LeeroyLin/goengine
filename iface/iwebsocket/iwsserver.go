package iwebsocket

import iface "github.com/LeeroyLin/goengine/iface/inetwork"

type IWSServer interface {
	Start()
	Stop()

	// AddRouter 添加路由业务
	AddRouter(msgId uint32, router WSRouterHandler)
	// SetDefaultRouter 设置默认路由
	SetDefaultRouter(router WSRouterHandler)
	// GetConnMgr 获得ConnMgr
	GetConnMgr() IWSConnManager
	// GetDataPack 获得封包解包对象
	GetDataPack() iface.IDataPack
}
