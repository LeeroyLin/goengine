package iwebsocket

type IWSConnManager interface {
	Add(conn IWSConnection)                  // 添加连接
	RemoveConn(conn IWSConnection)           // 删除连接
	Get(connID uint32) (IWSConnection, bool) // 利用ConnID获取连接
	Count() int                              // 获得当前连接数（精准，但效率很低）
	WeakCount() int                          // 获得当前连接数（非精准）
	Remove(connId uint32)                    // 移除连接
	StopAllConn()                            // 停止所有连接
}
