package network

import (
	"errors"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface/inetwork"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	// 隶属于的Server
	Server inetwork.IServer
	// 当前连接的socket TCP套接字
	conn *net.TCPConn
	// 当前连接的ID ID全局唯一
	connID uint32
	// 关闭通道
	closeChan chan interface{}
	// 是否关闭
	isClosed bool

	// 消息处理对象
	MsgHandler inetwork.IMsgHandler

	// 用于读数据和写数据两个Goroutine之间传递消息的通道 带缓冲
	msgBuffChan chan []byte

	// 属性Map
	property map[string]interface{}
	// 处理属性Map的读写锁
	propertyLock sync.RWMutex

	// 工作池数量
	workerPoolSize uint32
	// 最大消息队列通道容量
	maxMsgBuffChanLen uint32

	sync.RWMutex
}

// StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	elog.Info("[Conn] reader started.", c.connID)
	defer c.Stop()
	defer elog.Info("[Conn] reader exited.", c.connID)

	for {
		select {
		case <-c.closeChan:
			return
		default:
			dp := c.Server.GetDataPack()

			// 读取客户端的msg head
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
				elog.Error("[Conn] read msg head err:", err)
				return
			}

			// 拆包，得到 MsgId 和 DataLen
			msg, err := dp.Unpack(headData)
			if err != nil {
				elog.Error("[Conn] unpack msg err:", err)
				return
			}

			// 根据 DataLen 读取数据
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
					elog.Error("[Conn] read msg data err:", err)
					return
				}
			}
			msg.SetData(data)

			// 创建请求
			req := NewRequest(c, msg)

			// 有工作池
			if c.workerPoolSize > 0 {
				// 通过消息队列给工作池处理
				c.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				// 直接处理
				go c.MsgHandler.DoMsgHandler(req)
			}
		}
	}
}

// StartWriter 开启写消息Goroutine
func (c *Connection) StartWriter() {
	elog.Info("[Conn] writer started.", c.connID)
	defer elog.Info("[Conn] writer exited.", c.connID)

	for {
		select {
		case <-c.closeChan:
			return
		case data, ok := <-c.msgBuffChan:
			if ok {
				_, err := c.conn.Write(data)
				if err != nil {
					elog.Error("Write msg buff chan data err:", err)
					return
				}
			} else {
				elog.Error("msg buff chan has been closed.")
				return
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.Lock()
	c.isClosed = false
	c.closeChan = make(chan interface{})
	c.Unlock()

	go func() {
		select {
		case <-c.closeChan:
			return
		default:
			// 开启读业务
			go c.StartReader()
		}

		select {
		case <-c.closeChan:
			c.finalizer()
			return
		}
	}()
}

// Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	select {
	case <-c.closeChan:
		return
	default:
		close(c.closeChan)
		c.Lock()
		c.isClosed = true
		c.Unlock()
	}
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.connID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SendMsg 发送数据给客户端
func (c *Connection) SendMsg(msgId, serialId uint32, errCode uint16, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	// 封包
	dp := c.Server.GetDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, serialId, errCode, data))
	if err != nil {
		elog.Error("[Conn] pack msg err:", err)
		return err
	}

	_, err = c.GetTCPConnection().Write(msg)
	if err != nil {
		elog.Error("[Conn] send msg err:", err)
		return err
	}

	return nil
}

// SendBuffMsg 发送数据到客户端（带缓冲）
func (c *Connection) SendBuffMsg(msgId, serialId uint32, errCode uint16, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, c.maxMsgBuffChanLen)
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	// 封包
	dp := c.Server.GetDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, serialId, errCode, data))
	if err != nil {
		elog.Error("[Conn] pack msg err:", err)
		return err
	}

	select {
	case c.msgBuffChan <- msg:
		return nil
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	}
}

// GetProperty 获得属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	v, ok := c.property[key]
	if ok {
		return v, nil
	}

	return nil, errors.New("[Conn] property not found: " + key)
}

// SetProperty 设置属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// RemoveProperty 移除属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) finalizer() {
	// 从连接管理中移除该连接
	c.Server.GetConnMgr().RemoveConn(c)

	elog.Info("[Conn] connection stoped.", c.connID)

	if c.msgBuffChan == nil {
		return
	}

	// 关闭该连接全部管道
	select {
	case <-c.msgBuffChan:
		return
	default:
		close(c.msgBuffChan)
	}
}

// NewConnection 创建连接的方法
func NewConnection(workerPoolSize, maxMsgBuffChanLen uint32, server inetwork.IServer, conn *net.TCPConn, connID uint32, msgHandler inetwork.IMsgHandler) inetwork.IConnection {
	c := &Connection{
		Server:            server,
		conn:              conn,
		connID:            connID,
		closeChan:         make(chan interface{}),
		MsgHandler:        msgHandler,
		msgBuffChan:       nil,
		property:          nil,
		isClosed:          false,
		workerPoolSize:    workerPoolSize,
		maxMsgBuffChanLen: maxMsgBuffChanLen,
	}

	c.Server.GetConnMgr().Add(c)

	return c
}
