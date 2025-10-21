package msgcenter

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/def"
	"sync"
)

// MsgCenter 消息中心
type MsgCenter struct {
	title    string
	handlers map[uint32]def.MsgHandler
	msgChan  chan *def.MsgReqBundle
	sync.RWMutex
	closeChan chan interface{}
}

// NewMsgCenter 新建消息中心
//
// parameter msgCapacity 消息队列容量
func NewMsgCenter(title string, msgCapacity int, closeChan chan interface{}) *MsgCenter {
	mc := &MsgCenter{
		title:     title,
		handlers:  make(map[uint32]def.MsgHandler),
		msgChan:   make(chan *def.MsgReqBundle, msgCapacity),
		closeChan: closeChan,
	}

	return mc
}

// AddHandler 添加消息处理函数
func (center *MsgCenter) AddHandler(commId uint32, handler def.MsgHandler) {
	center.Lock()
	defer center.Unlock()

	_, ok := center.handlers[commId]
	if ok {
		elog.Error("[MsgCenter] already has same CommId handler. replaced.", center.title, commId)
	}
	center.handlers[commId] = handler
}

// RemoveHandler 删除对应通信id的消息处理函数
func (center *MsgCenter) RemoveHandler(commId uint32) {
	center.Lock()
	defer center.Unlock()

	delete(center.handlers, commId)
}

// ClearHandlers 清空所有消息处理函数
func (center *MsgCenter) ClearHandlers() {
	center.Lock()
	defer center.Unlock()

	for opId, _ := range center.handlers {
		delete(center.handlers, opId)
	}
}

// CloseMsgChan 关闭消息通道
func (center *MsgCenter) CloseMsgChan() {
	close(center.msgChan)
}

// Close 关闭消息中心
func (center *MsgCenter) Close() {
	center.ClearHandlers()
	center.CloseMsgChan()
}

// Run 运行消息中心，处理消息
func (center *MsgCenter) Run() {
	go func() {
		for {
			select {
			case <-center.closeChan:
				elog.Info("[MsgCenter] msg center closed.", center.title)
				return
			case bundle := <-center.msgChan:
				// 执行异步调用
				center.callAsync(bundle)
			}
		}
	}()
}

// Call 同步调用
func (center *MsgCenter) Call(mReq def.ICommReq) (interface{}, error) {
	center.RLock()
	defer center.RUnlock()

	commId := mReq.GetCommId()

	h, ok := center.handlers[commId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[MsgCenter] can not call msg. can not find msg id. center:%v CommId:%v", center.title, commId))
	}

	return h(true, mReq)
}

// CallAsync 异步调用 如果消息队列满了会阻塞
func (center *MsgCenter) CallAsync(mReq def.ICommReq, cb func(resp interface{}, err error)) {
	// 放入消息通道
	center.msgChan <- &def.MsgReqBundle{
		Req: mReq,
		Cb:  cb,
	}
}

// Cast 消息投递 忽略返回信息和错误信息
func (center *MsgCenter) Cast(mReq def.ICommReq) {
	commId := mReq.GetCommId()

	bundle := &def.MsgReqBundle{
		Req: mReq,
		Cb:  nil,
	}

	select {
	case center.msgChan <- bundle:
		break
	default:
		// 投递失败
		elog.Errorf("[MsgCenter] can not call msg. chan is full. msg dropped. center:%v CommId:%v\n", center.title, commId)
		break
	}
}

// 执行异步调用
func (center *MsgCenter) callAsync(bundle *def.MsgReqBundle) {
	commId := bundle.Req.GetCommId()

	center.RLock()

	h, ok := center.handlers[commId]

	center.RUnlock()

	if !ok {
		if bundle.Cb == nil {
			return
		}

		err := errors.New(fmt.Sprintf("[MsgCenter] can not call msg async. can not find msg id. center:%v CommId:%v", center.title, commId))
		bundle.Cb(nil, err)
		return
	}

	// 异步回调
	go func() {
		isSync := bundle.Cb != nil

		resp, err := h(isSync, bundle.Req)
		if err != nil {
			elog.Error("[MsgCenter] call handler err: ", center.title, commId, err)
			return
		}

		if bundle.Cb != nil {
			bundle.Cb(resp, nil)
		}
	}()
}
