package iface

type IModule interface {
	SetLife(life IModuleLife)
	GetName() string            // 获得模块名
	GetMsgCenter() IMsgCenter   // 获得消息中心
	GetDispatcher() IDispatcher // 获得模块间消息分发器
	DoInit(dispatcher IDispatcher, rpcGetter IRPCGetter, etcdGetter IETCDGetter, msgChanCapacity int, closeChan chan interface{})
	DoRun() error
	DoStop() error
	DoBeforeStop() error
}
