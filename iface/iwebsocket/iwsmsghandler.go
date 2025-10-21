package iwebsocket

type IWSMsgHandler interface {
	// DoMsgHandler 以非阻塞形式处理消息
	DoMsgHandler(request IWSRequest)
	// AddRouter 添加路由
	AddRouter(msgId uint32, router WSRouterHandler)
	// SetDefaultRouter 设置默认路由
	SetDefaultRouter(router WSRouterHandler)
	// StartWorkerPool 开启工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker处理
	SendMsgToTaskQueue(request IWSRequest)
}
