package network

import (
	"fmt"
	"github.com/LeeroyLin/goengine/core/syncmap"
	"github.com/LeeroyLin/goengine/iface/inetwork"
)

type ConnManager struct {
	connections *syncmap.SyncMap[uint32, inetwork.IConnection]
}

// Add 添加连接
func (cm *ConnManager) Add(conn inetwork.IConnection) {
	cm.connections.Add(conn.GetConnID(), conn)
}

// RemoveConn 删除连接
func (cm *ConnManager) RemoveConn(conn inetwork.IConnection) {
	cm.Remove(conn.GetConnID())
}

// Get 利用ConnId获得连接
func (cm *ConnManager) Get(connId uint32) (inetwork.IConnection, bool) {
	return cm.connections.Get(connId)
}

// Size 获得当前连接数
func (cm *ConnManager) Size() int {
	return cm.connections.Size()
}

// Remove 移除连接
func (cm *ConnManager) Remove(connId uint32) {
	conn, ok := cm.connections.GetAndDelete(connId)
	if ok {
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
	cm.connections.Range(func(connId uint32, conn inetwork.IConnection) bool {
		conn.(inetwork.IConnection).Stop()
		return true
	})
}

func NewConnManager() inetwork.IConnManager {
	cm := &ConnManager{
		connections: syncmap.NewSyncMap[uint32, inetwork.IConnection](),
	}

	return cm
}
