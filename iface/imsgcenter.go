package iface

import "github.com/LeeroyLin/goengine/def"

type IMsgCenter interface {
	// AddHandler 添加消息处理函数
	AddHandler(commId uint32, handler def.MsgHandler)

	// RemoveHandler 删除对应通信id的消息处理函数
	RemoveHandler(commId uint32)

	// ClearHandlers 清空所有消息处理函数
	ClearHandlers()

	// CloseMsgChan 关闭消息通道
	CloseMsgChan()

	// Close 关闭消息中心
	Close()

	// Run 运行消息中心，处理消息
	Run()

	// Call 同步调用
	Call(mReq def.ICommReq) (interface{}, error)

	// CallAsync 异步调用 如果消息队列满了会阻塞
	CallAsync(mReq def.ICommReq, cb func(resp interface{}, err error))

	// Cast 消息投递 忽略返回信息和错误信息
	Cast(mReq def.ICommReq)
}
