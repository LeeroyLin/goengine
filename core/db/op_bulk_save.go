package db

import (
	"context"
	"errors"
	"github.com/LeeroyLin/goengine/core/elog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
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

	buf := DBBufferPool.Get()
	vm, err := bsonrw.NewBSONValueWriter(buf)

	if err != nil {
		return nil, err
	}

	enc, err := bson.NewEncoder(vm)

	if err != nil {
		return nil, err
	}

	for i := 0; i < cnt; i++ {
		each := op.OpEachArr[i]

		buf.Reset()
		err = enc.Encode(each.Data)
		if err != nil {
			t := i
			elog.Error("[MongoDB] bulk save encode err.", op.DBName, op.CollName, t, err)
			continue
		}
		bytesCnt += buf.Len()

		elog.Debug("[MongoDB] bulk save", i, buf.Len())

		wm := mongo.NewReplaceOneModel().SetFilter(each.Filter).SetReplacement(buf.Bytes()).SetUpsert(true)
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
		_, err = c.BulkWrite(ctx, writeModels[startIdx:endIdx])
		cancel()

		elog.Debug("[MongoDB] bulk save split.", op.DBName, op.CollName, startIdx, endIdx)

		if err != nil {
			elog.Error("[MongoDB] bulk save err.", op.DBName, op.CollName, err)
			return nil, err
		}

		startIdx = endIdx
	}

	return nil, nil
}
