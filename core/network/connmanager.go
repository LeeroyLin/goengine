package network

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/iface/inetwork"
	"strconv"
	"sync"
)

type ConnManager struct {
	connections sync.Map
}

// Add 添加连接
func (cm *ConnManager) Add(conn inetwork.IConnection) {
	cm.connections.Store(conn.GetConnID(), conn)
}

// RemoveConn 删除连接
func (cm *ConnManager) RemoveConn(conn inetwork.IConnection) {
	cm.Remove(conn.GetConnID())
}

// Get 利用ConnId获得连接
func (cm *ConnManager) Get(connId uint32) (inetwork.IConnection, error) {
	conn, ok := cm.connections.Load(connId)
	if !ok {
		return nil, errors.New("[ConnMgr] connection not found. ConnId:" + strconv.Itoa(int(connId)))
	}

	return conn.(inetwork.IConnection), nil
}

// Len 获得当前连接数
func (cm *ConnManager) Len() int {
	cnt := 0

	cm.connections.Range(func(key any, value any) bool {
		cnt++

		return true
	})

	return cnt
}

// Remove 移除连接
func (cm *ConnManager) Remove(connId uint32) {
	conn, ok := cm.connections.Load(connId)
	if ok {
		cm.connections.Delete(connId)
		err := conn.(inetwork.IConnection).GetTCPConnection().Close()
		if err != nil {
			fmt.Println("[ConnMgr] close conn error ", err)
			return
		}
	}
}

// StopAllConn 停止所有连接
func (cm *ConnManager) StopAllConn() {
	fmt.Println("[ConnMgr] try stop all connections...")
	cm.connections.Range(func(key any, value any) bool {
		value.(inetwork.IConnection).Stop()
		return true
	})
}

func NewConnManager() inetwork.IConnManager {
	cm := &ConnManager{}

	return cm
}
