package inetwork

type IMsgHandler interface {
	// DoMsgHandler 以非阻塞形式处理消息
	DoMsgHandler(request IRequest)
	// AddRouter 添加路由
	AddRouter(msgId uint32, router RouterHandler)
	// SetDefaultRouter 设置默认路由
	SetDefaultRouter(router RouterHandler)
	// StartWorkerPool 开启工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker处理
	SendMsgToTaskQueue(request IRequest)
}
