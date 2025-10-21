package network

import (
	"github.com/LeeroyLin/goengine/iface/inetwork"
)

type Request struct {
	conn inetwork.IConnection // 已经和客户端建立好的连接
	msg  inetwork.IMessage    // 客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() inetwork.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求消息的id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

// NewRequest 新建Request结构
func NewRequest(conn inetwork.IConnection, msg inetwork.IMessage) inetwork.IRequest {
	r := &Request{
		conn: conn,
		msg:  msg,
	}

	return r
}
