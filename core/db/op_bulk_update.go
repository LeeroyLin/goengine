package db

import (
	"context"
	"errors"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBBulkUpdateOp 数据库批量更新操作
type DBBulkUpdateOp struct {
	DBOpBase
	OpEachArr []*DBOpEach
}

func NewDBBulkUpdateOp(fromModule, dbName, collName string) *DBBulkUpdateOp {
	op := &DBBulkUpdateOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		OpEachArr: make([]*DBOpEach, 0),
	}

	return op
}

func (op DBBulkUpdateOp) Exec(c *mongo.Collection) (interface{}, error) {
	cnt := len(op.OpEachArr)

	if cnt == 0 {
		elog.Info("[MongoDB] exec 0 length bulk update op.", op.DBName, op.CollName)
		return nil, errors.New("exec 0 length bulk update op")
	}

	writeModels := make([]mongo.WriteModel, 0)

	var endIdxArr []int
	bytesCnt := 0

	buf := DBBufferPool.Get()
	vm, err := bsonrw.NewBSONValueWriter(buf)

	if err != nil {
		elog.Error("[MongoDB] bulk update NewBSONValueWriter err.", op.DBName, op.CollName, err)
		return nil, err
	}

	enc, err := bson.NewEncoder(vm)

	if err != nil {
		elog.Error("[MongoDB] bulk update NewEncoder err.", op.DBName, op.CollName, err)
		return nil, err
	}

	for i := 0; i < cnt; i++ {
		each := op.OpEachArr[i]

		err = enc.Encode(each.Data)
		if err != nil {
			t := i
			elog.Error("[MongoDB] bulk update encode err.", op.DBName, op.CollName, t, err)
			continue
		}
		bytesCnt += buf.Len()

		wm := mongo.NewUpdateOneModel().SetFilter(each.Filter).SetUpdate(each.Data).SetUpsert(true)
		writeModels = append(writeModels, wm)

		if bytesCnt >= DB_Op_BulkWriteSize {
			endIdxArr = append(endIdxArr, i)
			bytesCnt = 0
		}
	}

	endIdxArr = append(endIdxArr, cnt)

	var startIdx = 0

	elog.Info("[MongoDB] bulk update split cnt.", op.DBName, op.CollName, len(endIdxArr))

	for _, endIdx := range endIdxArr {
		ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
		_, err = c.BulkWrite(ctx, writeModels[startIdx:endIdx])
		cancel()

		elog.Debug("[MongoDB] bulk update split.", op.DBName, op.CollName, startIdx, endIdx)

		if err != nil {
			elog.Error("[MongoDB] bulk update err.", op.DBName, op.CollName, err)
			return nil, err
		}

		startIdx = endIdx
	}

	return nil, nil
}
