package db

import (
	"context"
	"engine/core/elog"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBBulkDeleteOp 数据库批量删除操作
type DBBulkDeleteOp struct {
	DBOpBase
	OpEachArr []*DBOpEach
}

func NewDBBulkDeleteOp(fromModule, dbName, collName string) *DBBulkDeleteOp {
	op := &DBBulkDeleteOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
	}

	return op
}

func (op DBBulkDeleteOp) Exec(c *mongo.Collection) (interface{}, error) {
	cnt := len(op.OpEachArr)

	writeModels := make([]mongo.WriteModel, 0)

	for i := 0; i < cnt; i++ {
		each := op.OpEachArr[i]

		wm := mongo.NewDeleteOneModel().SetFilter(each.Filter)
		writeModels = append(writeModels, wm)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	_, err := c.BulkWrite(ctx, writeModels)
	cancel()

	if err != nil {
		elog.Error("[MongoDB] bulk delete err.", op.DBName, op.CollName, err)
		return nil, err
	}

	return nil, nil
}
