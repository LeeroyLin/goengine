package db

import (
	"context"
	"engine/core/elog"
	"engine/iface/idb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBWorker struct {
	url         string
	closeChan   chan interface{}
	mongoClient *mongo.Client
	rp          *readpref.ReadPref
	opChan      chan idb.IDBOp
	cbHandler   func(idb.IDBOp, interface{}, error)
}

func NewDBWorker(url string, rp *readpref.ReadPref) idb.IDBWorker {
	w := &DBWorker{
		url:       url,
		closeChan: make(chan interface{}),
		rp:        rp,
		opChan:    make(chan idb.IDBOp, 1024),
	}

	return w
}

func (w *DBWorker) Run() {
	if w.cbHandler == nil {
		elog.Error("[MongoDB] cbHandler have not been set yet.", w.url)
		return
	}

	clientOpts := options.Client().ApplyURI("mongodb://" + w.url).SetReadPreference(w.rp)
	clientOpts.SetConnectTimeout(DB_Conn_Timeout)
	clientOpts.SetSocketTimeout(DB_Conn_Timeout)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		elog.Error("[MongoDB] Connect db err.", w.url, err)
		return
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		elog.Error("[MongoDB] Connect db ping err.", w.url, err)
		return
	}

	w.mongoClient = client

	elog.Info("[MongoDB] Connect db success.", w.url)

	go w.exec()
}

func (w *DBWorker) Stop() {
	select {
	case <-w.closeChan:
		return
	default:
		close(w.closeChan)

		elog.Info("[MongoDB] Disconnect db.", w.url)

		err := w.mongoClient.Disconnect(context.Background())
		if err != nil {
			elog.Error("[MongoDB] Disconnect db err.", w.url, err)
			return
		}
	}
}

func (w *DBWorker) SetCBHandler(handler func(idb.IDBOp, interface{}, error)) {
	w.cbHandler = handler
}

func (w *DBWorker) Call(op idb.IDBOp) (interface{}, error) {
	coll := w.mongoClient.Database(op.GetDBName()).Collection(op.GetCollName())
	resData, err := op.Exec(coll)
	return resData, err
}

func (w *DBWorker) CastOp(dbOp idb.IDBOp) {
	w.opChan <- dbOp
}

func (w *DBWorker) exec() {
	for {
		select {
		case <-w.closeChan:
			return
		case op := <-w.opChan:
			coll := w.mongoClient.Database(op.GetDBName()).Collection(op.GetCollName())
			resData, err := op.Exec(coll)

			w.cbHandler(op, resData, err)
		}
	}
}
