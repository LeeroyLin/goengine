package db

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBBulkCreateIndexOp 数据库批量创建索引
type DBBulkCreateIndexOp struct {
	DBOpBase
	idxModels []mongo.IndexModel
}

func NewDBBulkCreateIndexOp(fromModule, dbName, collName string, idxModels []mongo.IndexModel) *DBBulkCreateIndexOp {
	op := &DBBulkCreateIndexOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		idxModels: idxModels,
	}

	return op
}

func (op DBBulkCreateIndexOp) Exec(c *mongo.Collection) (interface{}, error) {
	view := c.Indexes()

	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	_, err := view.CreateMany(ctx, op.idxModels)
	cancel()

	if err != nil {
		return nil, err
	}

	return nil, nil
}
