package def

// MsgHandler 处理消息的函数
type MsgHandler func(isSync bool, mr ICommReq) (interface{}, error)

// MsgRespHandler 异步回调函数
type MsgRespHandler func(resp interface{}, err error)

type ICommReq interface {
	GetCommId() uint32
}

type CommReqBase struct {
	CommId uint32
}

func (r *CommReqBase) GetCommId() uint32 {
	return r.CommId
}

// MsgReqBundle 消息请求包
type MsgReqBundle struct {
	Req ICommReq
	Cb  MsgRespHandler
}
