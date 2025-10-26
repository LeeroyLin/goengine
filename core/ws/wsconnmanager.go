package ws

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/syncmap"
	"github.com/LeeroyLin/goengine/iface/iwebsocket"
	"strconv"
)

type WSConnManager struct {
	connections *syncmap.SyncMap[uint32, iwebsocket.IWSConnection]
}

// Add 添加连接
func (cm *WSConnManager) Add(conn iwebsocket.IWSConnection) {
	cm.connections.Add(conn.GetConnID(), conn)
}

// RemoveConn 删除连接
func (cm *WSConnManager) RemoveConn(conn iwebsocket.IWSConnection) {
	cm.Remove(conn.GetConnID())
}

// Get 利用ConnId获得连接
func (cm *WSConnManager) Get(connId uint32) (iwebsocket.IWSConnection, bool) {
	return cm.connections.Get(connId)
}

// Size 获得当前连接数
func (cm *WSConnManager) Size() int {
	return cm.connections.Size()
}

// Remove 移除连接
func (cm *WSConnManager) Remove(connId uint32) {
	conn, ok := cm.connections.GetAndDelete(connId)
	if ok {
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
	cm.connections.Range(func(connId uint32, conn iwebsocket.IWSConnection) bool {
		conn.(iwebsocket.IWSConnection).Stop()
		return true
	})
}

func NewWSConnManager() iwebsocket.IWSConnManager {
	cm := &WSConnManager{
		connections: syncmap.NewSyncMap[uint32, iwebsocket.IWSConnection](),
	}

	return cm
}
