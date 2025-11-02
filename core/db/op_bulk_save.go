package db

import (
	"context"
	"errors"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/pool"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBBulkSaveOp 数据库批量存储操作
type DBBulkSaveOp struct {
	DBOpBase
	OpEachArr []*DBOpEach
}

func NewDBBulkSaveOp(fromModule, dbName, collName string) *DBBulkSaveOp {
	op := &DBBulkSaveOp{
		DBOpBase: DBOpBase{
			FromModule: fromModule,
			DBName:     dbName,
			CollName:   collName,
		},
		OpEachArr: make([]*DBOpEach, 0),
	}

	return op
}

func (op DBBulkSaveOp) Exec(c *mongo.Collection) (interface{}, error) {
	cnt := len(op.OpEachArr)

	if cnt == 0 {
		elog.Info("[MongoDB] exec 0 length bulk save op.", op.DBName, op.CollName)
		return nil, errors.New("exec 0 length bulk save op")
	}

	writeModels := make([]mongo.WriteModel, 0)

	var endIdxArr []int
	bytesCnt := 0

	for i := 0; i < cnt; i++ {
		each := op.OpEachArr[i]

		bytesCnt += each.DataBuffer.Len()

		wm := mongo.NewReplaceOneModel().SetFilter(each.Filter).
			SetReplacement(each.DataBuffer.AvailableBytes()).SetUpsert(true)
		writeModels = append(writeModels, wm)

		if bytesCnt >= DB_Op_BulkWriteSize {
			endIdxArr = append(endIdxArr, i)
			bytesCnt = 0
		}
	}

	endIdxArr = append(endIdxArr, cnt)

	var startIdx = 0

	elog.Info("[MongoDB] bulk save split cnt.", op.DBName, op.CollName, len(endIdxArr))

	for _, endIdx := range endIdxArr {
		ctx, cancel := context.WithTimeout(context.Background(), DB_Op_Timeout)
		_, err := c.BulkWrite(ctx, writeModels[startIdx:endIdx])
		cancel()

		elog.Debug("[MongoDB] bulk save split.", op.DBName, op.CollName, startIdx, endIdx)

		if err != nil {
			elog.Error("[MongoDB] bulk save err.", op.DBName, op.CollName, err)
			return nil, err
		}

		startIdx = endIdx
	}

	// 回收
	for _, each := range op.OpEachArr {
		pool.PFMBufferCtl.Set(each.DataBuffer)
		each.DataBuffer = nil
	}

	return nil, nil
}
