package db

import (
	"context"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBCreateIndexOp 数据库创建索引
type DBCreateIndexOp struct {
	DBOpBase
	idxModel mongo.IndexModel
}

func NewDBCreateIndexOp(fromModule, dbName, collName string, idxModel mongo.IndexModel) *DBCreateIndexOp {
	op := &DBCreateIndexOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		idxModel: idxModel,
	}

	return op
}

func (op DBCreateIndexOp) Exec(c *mongo.Collection) (interface{}, error) {
	view := c.Indexes()

	ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
	_, err := view.CreateOne(ctx, op.idxModel)
	cancel()

	if err != nil {
		elog.Error("[MongoDB] create index err.", op.DBName, op.CollName, err)
		return nil, err
	}

	return nil, nil
}
