package timer

import (
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"reflect"
	"runtime"
)

type DelayFunc struct {
	f    func(...interface{})
	args []interface{}
}

func NewDelayFunc(f func(...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

func (df *DelayFunc) String() string {
	pc := reflect.ValueOf(df.f).Pointer()
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("{DelayFunc:%s, args:%v}", funcName, df.args)
}

func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			elog.Error(df.String(), "Call err:", err)
		}
	}()

	df.f(df.args...)
}
