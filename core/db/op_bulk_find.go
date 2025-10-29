package db

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBBulkFindOp 数据库查找操作
type DBBulkFindOp struct {
	DBOpBase
	Filter     bson.M
	Options    []*options.FindOptions
	ObjCreator func() interface{}
}

func NewDBBulkFindOp(fromModule, dbName, collName string, filter bson.M, objCreator func() interface{}, opts ...*options.FindOptions) *DBBulkFindOp {
	op := &DBBulkFindOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Filter:     filter,
		Options:    nil,
		ObjCreator: objCreator,
	}

	if opts != nil {
		op.Options = opts
	}

	return op
}

func (op DBBulkFindOp) Exec(c *mongo.Collection) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)

	cursor, err := c.Find(ctx, op.Filter, op.Options...)
	defer cancel()

	if err != nil {
		return nil, err
	}

	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			elog.Error("[MongoDB] find close cursor err.", op.DBName, op.CollName, err)
		}
	}()

	objs := make([]interface{}, 0)

	for cursor.Next(ctx) {
		if err = ctx.Err(); err != nil {
			return nil, err
		}

		if err = cursor.Err(); err != nil {
			return nil, err
		}

		obj := op.ObjCreator()
		err = cursor.Decode(obj)
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}
