package elog

import (
	"fmt"
	"testing"
)

func TestStdLog(t *testing.T) {
	Debug("debug content1")
	Debug("debug content2")

	Debugf("debug a = %d\n", 123)

	SetLogFlags(BitDate | BitLongFile | BitLevel)
	Info("log info content")

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("log error content")

	AddLogFlag(BitShortFile | BitTime)
	Stack("log stack")

	SetLogFile("./log", "test1")
	Debug("===> log debug content ~~666")
	Debug("===> log debug content ~~888")
	Error("===> log Error!!!! ~~~555~~~")

	CloseDebug()
	Debug("===> should not show!!!")
	Debug("===> should not show!!!")
	Error("===> Error after debug close!!!")

	fmt.Println("test finished.")
}
