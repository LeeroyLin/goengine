package ws

import (
	"engine/iface/inetwork"
	"engine/iface/iwebsocket"
)

type WSRequest struct {
	conn iwebsocket.IWSConnection // 已经和客户端建立好的连接
	msg  inetwork.IMessage        // 客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *WSRequest) GetConnection() iwebsocket.IWSConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *WSRequest) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId 获取请求消息的id
func (r *WSRequest) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

// NewWSRequest 新建Request结构
func NewWSRequest(conn iwebsocket.IWSConnection, msg inetwork.IMessage) iwebsocket.IWSRequest {
	r := &WSRequest{
		conn: conn,
		msg:  msg,
	}

	return r
}
