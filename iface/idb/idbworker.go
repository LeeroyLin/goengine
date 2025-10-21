package idb

type IDBWorker interface {
	Run()
	Stop()
	Call(op IDBOp) (interface{}, error)
	CastOp(dbOp IDBOp)
	SetCBHandler(handler func(IDBOp, interface{}, error))
}
