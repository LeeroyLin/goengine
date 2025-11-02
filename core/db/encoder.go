package db

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/pool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

var DBBsonEncoder *dbBsonEncoder

// 通用bson解析器
type dbBsonEncoder struct {
	buffer  *pool.PFMBuffer
	vm      bsonrw.ValueWriter
	encoder *bson.Encoder
}

func init() {
	DBBsonEncoder = &dbBsonEncoder{
		buffer: pool.NewPFMBuffer(2048),
	}

	vm, err := bsonrw.NewBSONValueWriter(DBBsonEncoder.buffer)

	if err != nil {
		elog.Error("[DBCommon] new bson value writer err:", err)
		return
	}

	DBBsonEncoder.vm = vm

	encoder, err := bson.NewEncoder(vm)

	if err != nil {
		elog.Error("[DBCommon] new bson encoder err:", err)
		return
	}

	DBBsonEncoder.encoder = encoder
}

// EncodeWithPool 将对象转为bson字节数组（池化）
func (e *dbBsonEncoder) EncodeWithPool(v interface{}) (*pool.PFMBuffer, error) {
	e.buffer.Reset()

	err := e.encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	// 获得encode的数据byte数组
	dataBytes, l := e.buffer.Bytes()

	// 从池子中取一个buffer
	newBuf := pool.PFMBufferCtl.Get(l)

	// 将数据存入新buffer
	err = newBuf.WriteUtil(dataBytes, l)
	if err != nil {
		return nil, err
	}

	return newBuf, nil
}
