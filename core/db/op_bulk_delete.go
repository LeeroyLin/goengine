package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBBulkDeleteOp 数据库批量删除操作
type DBBulkDeleteOp struct {
	DBOpBase
	Filters []interface{}
}

func NewDBBulkDeleteOp(fromModule, dbName, collName string) *DBBulkDeleteOp {
	op := &DBBulkDeleteOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		Filters: make([]interface{}, 0),
	}

	return op
}

func (op DBBulkDeleteOp) Exec(c *mongo.Collection) (interface{}, error) {
	cnt := len(op.Filters)

	writeModels := make([]mongo.WriteModel, 0)

	for i := 0; i < cnt; i++ {
		filter := op.Filters[i]

		wm := mongo.NewDeleteOneModel().SetFilter(filter)
		writeModels = append(writeModels, wm)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	_, err := c.BulkWrite(ctx, writeModels)
	cancel()

	if err != nil {
		return nil, err
	}

	return nil, nil
}
