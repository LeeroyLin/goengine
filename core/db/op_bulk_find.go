package db

import (
	"context"
	"engine/core/elog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBBulkFindOp 数据库查找操作
type DBBulkFindOp struct {
	DBOpBase
	Filter     bson.M
	Skip       int64
	Limit      int64
	Sort       bson.M
	ObjCreator func() interface{}
}

func NewDBBulkFindOp(fromModule, dbName, collName string, filter bson.M, objCreator func() interface{}) *DBBulkFindOp {
	op := &DBBulkFindOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Filter:     filter,
		Skip:       0,
		Limit:      0,
		Sort:       nil,
		ObjCreator: objCreator,
	}

	return op
}

func (op DBBulkFindOp) Exec(c *mongo.Collection) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)

	findOptions := &options.FindOptions{
		Skip:  &op.Skip,
		Limit: &op.Limit,
	}

	if op.Sort != nil {
		findOptions.Sort = op.Sort
	}

	cursor, err := c.Find(ctx, op.Filter, findOptions)
	defer cancel()

	if err != nil {
		elog.Error("[MongoDB] find err.", op.DBName, op.CollName, err)
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
			elog.Error("[MongoDB] bulk find get next err", op.DBName, op.CollName, err)
			return nil, err
		}

		if err = cursor.Err(); err != nil {
			elog.Error("[MongoDB] bulk find cursor err", op.DBName, op.CollName, err)
			return nil, err
		}

		obj := op.ObjCreator()
		err = cursor.Decode(obj)
		if err != nil {
			elog.Error("[MongoDB] bulk find decode err.", op.DBName, op.CollName, err)
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}
