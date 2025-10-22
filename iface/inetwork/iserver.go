package inetwork

type IServer interface {
	Start()
	Stop()

	// AddRouter 添加路由业务
	AddRouter(msgId uint32, router RouterHandler)
	// SetDefaultRouter 设置默认路由
	SetDefaultRouter(router RouterHandler)
	// GetConnMgr 获得ConnMgr
	GetConnMgr() IConnManager
	// GetDataPack 获得封包解包对象
	GetDataPack() IDataPack
	// RecycleId 回收连接id
	RecycleId(connId uint32)
}
