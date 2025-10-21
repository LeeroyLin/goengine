package db

import (
	"engine/core/pool"
	"time"
)

type DBOpEach struct {
	Filter interface{}
	Data   interface{}
}

var DBBufferPool = pool.NewBufferPool(2048, 256)

const (
	DB_Conn_Timeout     = 5 * time.Second
	DB_Op_Timeout       = 5 * time.Second
	DB_Op_BulkWriteSize = 2 * 1024 * 1024
)
