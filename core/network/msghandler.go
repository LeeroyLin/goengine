package network

import (
	"engine/core/elog"
	"engine/iface/inetwork"
	"strconv"
)

type MsgHandler struct {
	// 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]inetwork.RouterHandler
	// Worker池数量
	WorkerPoolSize uint32
	// Worder对应最大消息队列数
	MaxWorkerTaskLen uint32
	// 消息队列
	TaskQueue []chan inetwork.IRequest
	// 默认路由
	DefaultRouter inetwork.RouterHandler
}

func (mh *MsgHandler) DoMsgHandler(req inetwork.IRequest) {
	msgId := req.GetMsgId()

	// 获取MsgId对应路由
	router, ok := mh.Apis[msgId]

	// 不存在
	if !ok {
		if mh.DefaultRouter != nil {
			mh.DefaultRouter(req)
		} else {
			elog.Errorf("[MsgHandler] msgHandler for msgId %d not found\n", msgId)
		}

		return
	}

	// 处理消息
	if router != nil {
		router(req)
	}
}

func (mh *MsgHandler) AddRouter(msgId uint32, router inetwork.RouterHandler) {
	// 是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		elog.Panic("[MsgHandler] repeated api. MsgId: " + strconv.Itoa(int(msgId)))
	}

	// 绑定
	mh.Apis[msgId] = router
	elog.Debug("[MsgHandler] Add api msgId: " + strconv.Itoa(int(msgId)))
}

func (mh *MsgHandler) SetDefaultRouter(router inetwork.RouterHandler) {
	mh.DefaultRouter = router
}

// StartOneWorker 开启一个Worker
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan inetwork.IRequest) {
	elog.Debug("[Worker] start worker id: " + strconv.Itoa(workerID))
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 开启工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 依次开启Worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 给Worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan inetwork.IRequest, mh.MaxWorkerTaskLen)

		// 开启一个Worder
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker处理
func (mh *MsgHandler) SendMsgToTaskQueue(request inetwork.IRequest) {
	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	mh.TaskQueue[workerID] <- request
}

func NewMsgHandler(workerPoolSize, maxWorkerTaskLen uint32) inetwork.IMsgHandler {
	mh := &MsgHandler{
		Apis:             make(map[uint32]inetwork.RouterHandler),
		WorkerPoolSize:   workerPoolSize,
		MaxWorkerTaskLen: maxWorkerTaskLen,
		TaskQueue:        make([]chan inetwork.IRequest, workerPoolSize),
	}

	return mh
}
