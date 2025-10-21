package elog

import (
	"os"
)

var StdLog = NewLog(os.Stdout, "", BitDefault)

func GetLogFlags() int {
	return StdLog.GetLogFlags()
}

func SetLogFlags(flag int) {
	StdLog.SetLogFlags(flag)
}

func AddLogFlag(flag int) {
	StdLog.AddLogFlag(flag)
}

func SetPrefix(prefix string) {
	StdLog.SetPrefix(prefix)
}

func SetLogFile(fileDir, fileName string) {
	StdLog.SetLogFile(fileDir, fileName)
}

func OpenDebug() {
	StdLog.OpenDebug()
}

func CloseDebug() {
	StdLog.CloseDebug()
}

func Debug(v ...interface{}) {
	StdLog.Debug(v...)
}

func Info(v ...interface{}) {
	StdLog.Info(v...)
}

func Error(v ...interface{}) {
	StdLog.Error(v...)
}

func Panic(v ...interface{}) {
	StdLog.Panic(v...)
}

func Fatal(v ...interface{}) {
	StdLog.Fatal(v...)
}

func Debugf(format string, v ...interface{}) {
	StdLog.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	StdLog.Infof(format, v...)
}

func Errorf(format string, v ...interface{}) {
	StdLog.Errorf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	StdLog.Panicf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	StdLog.Fatalf(format, v...)
}

func Stack(v ...interface{}) {
	StdLog.Stack(v...)
}
