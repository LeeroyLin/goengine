package ws

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface/iwebsocket"
	"strconv"
)

type WSMsgHandler struct {
	// 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]iwebsocket.WSRouterHandler
	// Worker池数量
	WorkerPoolSize uint32
	// Worder对应最大消息队列数
	MaxWorkerTaskLen uint32
	// 消息队列
	TaskQueue []chan iwebsocket.IWSRequest
	// 默认路由
	DefaultRouter iwebsocket.WSRouterHandler
}

func (mh *WSMsgHandler) DoMsgHandler(req iwebsocket.IWSRequest) {
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

func (mh *WSMsgHandler) AddRouter(msgId uint32, router iwebsocket.WSRouterHandler) {
	// 是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		elog.Panic("[MsgHandler] repeated api. MsgId: " + strconv.Itoa(int(msgId)))
	}

	// 绑定
	mh.Apis[msgId] = router
	elog.Debug("[MsgHandler] Add api msgId: " + strconv.Itoa(int(msgId)))
}

func (mh *WSMsgHandler) SetDefaultRouter(router iwebsocket.WSRouterHandler) {
	mh.DefaultRouter = router
}

// StartOneWorker 开启一个Worker
func (mh *WSMsgHandler) StartOneWorker(workerID int, taskQueue chan iwebsocket.IWSRequest) {
	elog.Debug("[Worker] start worker id: " + strconv.Itoa(workerID))
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 开启工作池
func (mh *WSMsgHandler) StartWorkerPool() {
	// 依次开启Worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 给Worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan iwebsocket.IWSRequest, mh.MaxWorkerTaskLen)

		// 开启一个Worder
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker处理
func (mh *WSMsgHandler) SendMsgToTaskQueue(request iwebsocket.IWSRequest) {
	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	mh.TaskQueue[workerID] <- request
}

func NewWSMsgHandler(workerPoolSize, maxWorkerTaskLen uint32) iwebsocket.IWSMsgHandler {
	mh := &WSMsgHandler{
		Apis:             make(map[uint32]iwebsocket.WSRouterHandler),
		WorkerPoolSize:   workerPoolSize,
		MaxWorkerTaskLen: maxWorkerTaskLen,
		TaskQueue:        make([]chan iwebsocket.IWSRequest, workerPoolSize),
	}

	return mh
}
