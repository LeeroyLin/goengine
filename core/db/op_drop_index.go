package db

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBDropIndexOp 数据库删除索引
type DBDropIndexOp struct {
	DBOpBase
	Name string
	Key  interface{}
}

func NewDBDropIndexOp(fromModule, dbName, collName string, name string) *DBDropIndexOp {
	op := &DBDropIndexOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Name: name,
	}

	return op
}

func NewDBDropIndexOpWithKey(fromModule, dbName, collName string, key interface{}) *DBDropIndexOp {
	op := &DBDropIndexOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Key: key,
	}

	return op
}

func (op DBDropIndexOp) Exec(c *mongo.Collection) (interface{}, error) {
	view := c.Indexes()

	if op.Key != nil {
		ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
		_, err := view.DropOneWithKey(ctx, op.Key)
		cancel()

		if err != nil {
			elog.Error("[MongoDB] drop index with key err.", op.DBName, op.CollName, op.Key, err)
			return nil, err
		}

		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	_, err := view.DropOne(ctx, op.Name)
	cancel()

	if err != nil {
		elog.Error("[MongoDB] drop index with name err.", op.DBName, op.CollName, op.Name, err)
		return nil, err
	}

	return nil, nil
}
