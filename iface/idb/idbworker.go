package idb

type IDBWorker interface {
	Run() error
	Stop() error
	Call(op IDBOp) (interface{}, error)
	CastOp(dbOp IDBOp)
	SetCBHandler(handler func(IDBOp, interface{}, error))
}
