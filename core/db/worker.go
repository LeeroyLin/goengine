package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface/idb"
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

func (w *DBWorker) Run() error {
	if w.cbHandler == nil {
		return errors.New(fmt.Sprint("[MongoDB] cbHandler have not been set yet", w.url))
	}

	clientOpts := options.Client().ApplyURI("mongodb://" + w.url).SetReadPreference(w.rp)
	clientOpts.SetConnectTimeout(DB_Conn_Timeout)
	clientOpts.SetSocketTimeout(DB_Conn_Timeout)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return errors.New(fmt.Sprint("[MongoDB] Connect db err.", w.url, err))
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return errors.New(fmt.Sprint("[MongoDB] Connect db ping err.", w.url, err))
	}

	w.mongoClient = client

	elog.Info("[MongoDB] Connect db success.", w.url)

	go w.exec()

	return nil
}

func (w *DBWorker) Stop() error {
	select {
	case <-w.closeChan:
		return nil
	default:
		close(w.closeChan)

		elog.Info("[MongoDB] Disconnect db.", w.url)

		err := w.mongoClient.Disconnect(context.Background())
		if err != nil {
			elog.Fatal("[MongoDB] Disconnect db err.", w.url, err)
			return err
		}
	}

	return nil
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
