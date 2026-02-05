package network

import (
	"github.com/LeeroyLin/goengine/iface/inetwork"
)

type Message struct {
	// 消息的ID
	Id uint32
	// 序号
	Serial uint32
	// 错误码
	ErrCode uint16
	// 消息的长度
	DataLen uint32
	// 消息的内容
	Data []byte
}

// NewMsgPackage 创建一个Message消息包
func NewMsgPackage(id, serial uint32, errCode uint16, data []byte) inetwork.IMessage {
	return &Message{
		Id:      id,
		Serial:  serial,
		ErrCode: errCode,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// GetMsgId 获取消息ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// GetSerialId 获取序号ID
func (msg *Message) GetSerialId() uint32 {
	return msg.Serial
}

// GetErrCode 获取错误码
func (msg *Message) GetErrCode() uint16 {
	return msg.ErrCode
}

// GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

// SetMsgId 设计消息ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

// SetSerialId 设置序号ID
func (msg *Message) SetSerialId(serial uint32) {
	msg.Serial = serial
}

// SetErrCode 设置错误码
func (msg *Message) SetErrCode(errCode uint16) {
	msg.ErrCode = errCode
}

// SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
