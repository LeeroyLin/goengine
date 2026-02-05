package inetwork

type IMessage interface {
	GetDataLen() uint32  // 获取消息数据段长度
	GetMsgId() uint32    // 获取消息ID
	GetSerialId() uint32 // 获取序号
	GetErrCode() uint16  // 获取错误码
	GetData() []byte     // 获取消息内容

	SetMsgId(uint32)    // 设置消息ID
	SetSerialId(uint32) // 设置序号
	SetErrCode(uint16)  // 设置错误码
	SetData([]byte)     // 设置消息内容
	SetDataLen(uint32)  // 设置消息数据段长度
}
