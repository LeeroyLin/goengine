package iface

import "github.com/LeeroyLin/goengine/def"

type IDispatcher interface {
	// Call 同步调用模块
	Call(module string, req def.ICommReq) (interface{}, error)

	// CallAsync 异步调用模块
	CallAsync(module string, req def.ICommReq, cb def.MsgRespHandler) error

	// Cast 向模块投递消息
	Cast(module string, req def.ICommReq)
}
