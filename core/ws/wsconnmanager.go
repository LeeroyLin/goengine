package ws

import (
	"engine/iface/iwebsocket"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type WSConnManager struct {
	connections sync.Map
}

// Add 添加连接
func (cm *WSConnManager) Add(conn iwebsocket.IWSConnection) {
	cm.connections.Store(conn.GetConnID(), conn)
}

// RemoveConn 删除连接
func (cm *WSConnManager) RemoveConn(conn iwebsocket.IWSConnection) {
	cm.Remove(conn.GetConnID())
}

// Get 利用ConnId获得连接
func (cm *WSConnManager) Get(connId uint32) (iwebsocket.IWSConnection, error) {
	conn, ok := cm.connections.Load(connId)
	if !ok {
		return nil, errors.New("[ConnMgr] connection not found. ConnId:" + strconv.Itoa(int(connId)))
	}

	return conn.(iwebsocket.IWSConnection), nil
}

// Len 获得当前连接数
func (cm *WSConnManager) Len() int {
	cnt := 0

	cm.connections.Range(func(key any, value any) bool {
		cnt++

		return true
	})

	return cnt
}

// Remove 移除连接
func (cm *WSConnManager) Remove(connId uint32) {
	conn, ok := cm.connections.Load(connId)
	if ok {
		cm.connections.Delete(connId)
		err := conn.(iwebsocket.IWSConnection).GetTCPConnection().Close()
		if err != nil {
			fmt.Println("[ConnMgr] close conn error ", err)
			return
		}
	}
}

// StopAllConn 停止所有连接
func (cm *WSConnManager) StopAllConn() {
	fmt.Println("[ConnMgr] try stop all connections...")
	cm.connections.Range(func(key any, value any) bool {
		value.(iwebsocket.IWSConnection).Stop()
		return true
	})
}

func NewWSConnManager() iwebsocket.IWSConnManager {
	cm := &WSConnManager{}

	return cm
}
