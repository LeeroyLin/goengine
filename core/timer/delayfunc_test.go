package timer

import (
	"fmt"
	"testing"
)

func PrintFunc(v ...interface{}) {
	fmt.Println(v...)
}

func TestDelayFunc(t *testing.T) {
	df := NewDelayFunc(PrintFunc, []interface{}{"1", "2", "3"})
	fmt.Println("df.String() = ", df.String())
	df.Call()
}
