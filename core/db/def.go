package db

import (
	"github.com/LeeroyLin/goengine/core/pool"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type DBOpEach struct {
	Filter     bson.M
	DataBuffer *pool.PFMBuffer
}

func WrapDBOpEach(filter bson.M, data interface{}) *DBOpEach {
	// 利用池化将对象转为bson字节buffer
	dataBuffer, err := DBBsonEncoder.EncodeWithPool(data)
	if err != nil {
		return nil
	}

	return &DBOpEach{
		Filter:     filter,
		DataBuffer: dataBuffer,
	}
}

const (
	DB_Conn_Timeout     = 5 * time.Second
	DB_Op_Timeout       = 5 * time.Second
	DB_Op_BulkWriteSize = 2 * 1024 * 1024
)
