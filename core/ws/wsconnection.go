package ws

import (
	"errors"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/network"
	"github.com/LeeroyLin/goengine/iface/iwebsocket"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"sync"
	"syscall"
	"time"
)

type WSConnection struct {
	// 隶属于的Server
	Server iwebsocket.IWSServer
	// 当前连接的socket
	conn *websocket.Conn
	// 当前连接的ID ID全局唯一
	connID uint32
	// 关闭通道
	closeChan chan interface{}
	// 是否关闭
	isClosed bool

	// 消息处理对象
	MsgHandler iwebsocket.IWSMsgHandler

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
func (c *WSConnection) StartReader() {
	elog.Info("[Conn] reader started.", c.connID)
	defer c.Stop()
	defer elog.Info("[Conn] reader exited.", c.connID)

	for {
		select {
		case <-c.closeChan:
			return
		default:
			dp := c.Server.GetDataPack()

			_, reader, err := c.conn.NextReader()
			if err != nil {
				if errors.Is(err, syscall.ECONNRESET) {
					return
				}
				if errors.Is(err, websocket.ErrCloseSent) {
					return
				}
				if errors.Is(err, io.EOF) {
					return
				}

				elog.Error("[Conn] read msg get reader err:", err)
				return
			}

			// 读取客户端的msg head
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(reader, headData); err != nil {
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
				if _, err := io.ReadFull(reader, data); err != nil {
					elog.Error("[Conn] read msg data err:", err)
					return
				}
			}
			msg.SetData(data)

			// 创建请求
			req := NewWSRequest(c, msg)

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
func (c *WSConnection) StartWriter() {
	elog.Info("[Conn] writer started.", c.connID)
	defer elog.Info("[Conn] writer exited.", c.connID)

	for {
		select {
		case <-c.closeChan:
			return
		case data, ok := <-c.msgBuffChan:
			if ok {
				err := c.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					elog.Error("Write msg buff chan data err:", c.connID, err)
					return
				}

			} else {
				elog.Error("msg buff chan has been closed.", c.connID)
				return
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *WSConnection) Start() {
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
func (c *WSConnection) Stop() {
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

// GetTCPConnection 从当前连接获取原始的socket
func (c *WSConnection) GetTCPConnection() *websocket.Conn {
	return c.conn
}

// GetConnID 获取当前连接ID
func (c *WSConnection) GetConnID() uint32 {
	return c.connID
}

// RemoteAddr 获取远程客户端地址信息
func (c *WSConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SendMsg 发送数据给客户端
func (c *WSConnection) SendMsg(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	// 封包
	dp := c.Server.GetDataPack()
	msg, err := dp.Pack(network.NewMsgPackage(msgId, data))
	if err != nil {
		elog.Error("[Conn] pack msg err:", c.connID, err)
		return err
	}

	w, err := c.GetTCPConnection().NextWriter(websocket.BinaryMessage)
	if err != nil {
		elog.Error("[Conn] get next writer err:", c.connID, err)
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		elog.Error("[Conn] send msg err:", c.connID, err)
		return err
	}

	err = w.Close()
	if err != nil {
		elog.Error("[Conn] close writer err:", c.connID, err)
		return err
	}

	return nil
}

// SendBuffMsg 发送数据到客户端（带缓冲）
func (c *WSConnection) SendBuffMsg(msgId uint32, data []byte) error {
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
	msg, err := dp.Pack(network.NewMsgPackage(msgId, data))
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
func (c *WSConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	v, ok := c.property[key]
	if ok {
		return v, nil
	}

	return nil, errors.New("[Conn] property not found: " + key)
}

// SetProperty 设置属性
func (c *WSConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// RemoveProperty 移除属性
func (c *WSConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *WSConnection) finalizer() {
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

// NewWSConnection 创建连接的方法
func NewWSConnection(workerPoolSize, maxMsgBuffChanLen uint32, server iwebsocket.IWSServer, conn *websocket.Conn, connID uint32, msgHandler iwebsocket.IWSMsgHandler) iwebsocket.IWSConnection {
	c := &WSConnection{
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
