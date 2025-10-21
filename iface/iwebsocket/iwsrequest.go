package iwebsocket

// IWSRequest 将连接信息和请求数据封装在Request
type IWSRequest interface {
	GetConnection() IWSConnection // 获取请求连接信息
	GetData() []byte              // 获取请求消息的数据
	GetMsgId() uint32             // 获取请求消息id
}
