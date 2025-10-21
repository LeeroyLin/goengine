package idb

import "go.mongodb.org/mongo-driver/mongo"

type IDBOp interface {
	GetDBName() string
	GetCollName() string
	GetFromModule() string
	Exec(c *mongo.Collection) (interface{}, error)
}
