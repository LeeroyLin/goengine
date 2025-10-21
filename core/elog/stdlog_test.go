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
	Info("zinx info content")

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("zinx error content")

	AddLogFlag(BitShortFile | BitTime)
	Stack("zinx stack")

	SetLogFile("./log", "zinx.log")
	Debug("===> zinx debug content ~~666")
	Debug("===> zinx debug content ~~888")
	Error("===> zinx Error!!!! ~~~555~~~")

	CloseDebug()
	Debug("===> should not show!!!")
	Debug("===> should not show!!!")
	Error("===> Error after debug close!!!")

	fmt.Println("test finished.")
}
