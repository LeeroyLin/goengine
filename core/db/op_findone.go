package db

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBFindOneOp 数据库查找一个数据操作
type DBFindOneOp struct {
	DBOpBase
	Filter     bson.M
	ObjCreator func() interface{}
}

func NewDBFindOneOp(fromModule, dbName, collName string, filter bson.M, objCreator func() interface{}) *DBFindOneOp {
	op := &DBFindOneOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Filter:     filter,
		ObjCreator: objCreator,
	}

	return op
}

func (op DBFindOneOp) Exec(c *mongo.Collection) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	res := c.FindOne(ctx, op.Filter)
	cancel()

	err := res.Err()
	if err != nil {
		elog.Error("[MongoDB] find one err.", op.DBName, op.CollName, err)
		return nil, err
	}

	obj := op.ObjCreator()
	err = res.Decode(obj)
	if err != nil {
		elog.Error("[MongoDB] find one decode err.", op.DBName, op.CollName, err)
		return nil, err
	}

	return obj, nil
}
