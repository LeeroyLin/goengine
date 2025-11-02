package test

import (
	"github.com/LeeroyLin/goengine/core/db"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface/idb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

type DocAccount struct {
	PlayerId    int64     `bson:"_id"`
	Username    string    `bson:"username"`
	Password    string    `bson:"password"`
	CreateAt    time.Time `bson:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt"`
}

func NewDocAccountFace() interface{} {
	return &DocAccount{}
}

func TestFindOne(t *testing.T) {
	//filter := bson.M{
	//	"_id": 7389222314677637120,
	//}
	filter := bson.M{
		"username": "1231234",
	}

	op := db.NewDBFindOneOp("test", "dev_account", "account", filter, NewDocAccountFace)

	runDBWorker(op)
}

func TestBulkFind(t *testing.T) {
	tm := time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC)

	filter := bson.M{
		"createdAt": bson.M{
			"$gt": tm,
		},
	}

	op := db.NewDBBulkFindOp("test", "dev_account", "account", filter, NewDocAccountFace)

	runDBWorker(op)
}

func TestBulkDelete(t *testing.T) {
	filters := []bson.M{
		{
			"username": "321321",
		},
		{
			"username": "1231233",
		},
	}

	op := db.NewDBBulkDeleteOp("test", "dev_account", "account", filters)

	runDBWorker(op)
}

func runDBWorker(op idb.IDBOp) {
	worker := db.NewDBWorker("admin:Mongo666,.@inner.archimetagame.com:27017", readpref.Primary())
	worker.SetCBHandler(func(op idb.IDBOp, i interface{}, err error) {

	})

	err := worker.Run()
	if err != nil {
		elog.Error("start worker err:", err)
		return
	}

	res, err := worker.Call(op)
	if err != nil {
		elog.Error("res err:", err)
		return
	}

	elog.Infof("res: %v", res)

	time.Sleep(time.Hour)
}
