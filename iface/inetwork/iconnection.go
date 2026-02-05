package inetwork

import (
	"net"
)

type IConnection interface {
	// Start 启动连接，让当前连接开始工作
	Start()
	// Stop 停止连接，结束当前连接状态
	Stop()
	// GetTCPConnection 从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// SendMsg 发送数据给客户端
	SendMsg(msgId, serialId uint32, errCode uint16, data []byte) error
	// SendBuffMsg 发送数据给客户端（带缓冲）
	SendBuffMsg(msgId, serialId uint32, errCode uint16, data []byte) error
	// GetProperty 获取属性
	GetProperty(name string) (interface{}, error)
	// SetProperty 设置属性
	SetProperty(name string, value interface{})
	// RemoveProperty 移除属性
	RemoveProperty(name string)
}

// HandFunc 定义一个统一处理连接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
