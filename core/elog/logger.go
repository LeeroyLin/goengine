package elog

import (
	"bytes"
	"engine/core/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	LOG_MAX_BUF = 1024 * 1024
)

// 日志头部信息标记位，采用bitmap方式，用户可以选择头部需要哪些标记位被打印
const (
	BitDate         = 1 << iota             // 日期标记位  2019/01/23
	BitTime                                 // 时间标记位  01:23:12
	BitMicroSeconds                         // 微秒级标记位 01:23:12.111222
	BitLongFile                             // 完整文件名称 /home/go/src/zinx/server.go
	BitShortFile                            // 最后文件名   server.go
	BitLevel                                // 当前日志级别： 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
	BitStdFlag      = BitDate | BitTime     // 标准头部日志格式
	BitDefault      = BitLevel | BitStdFlag // 默认日志头部格式
)

// 日志级别
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
	LogPanic
	LogFatal
)

// 日志级别的显示字符串
var levels = []string{
	"[DEBUG] ",
	"[INFO]  ",
	"[WARN]  ",
	"[ERROR] ",
	"[PANIC] ",
	"[FATAL] ",
}

type Logger struct {
	// 确保多协程读写文件，防止文件内容混乱，做到协程安全
	mu sync.Mutex
	// 每行log日志的前缀字符串,拥有日志标记
	prefix string
	// 日志标记位
	flag int
	// 日志输出的文件描述符
	out io.Writer
	// 输出的缓冲区
	buf bytes.Buffer
	// 当前日志绑定的输出文件
	file *os.File
	// 是否打印调试debug信息
	debug bool
	// 获取日志文件名和代码上述的runtime.Call 的函数调用层数
	callDepth int
}

// NewLog 创建一个日志
// Parameters:
//
//	out: 日志输出文件描述符
//	prefix: 日志前缀
//	flag: 日志标记位
//
// Return:
//
//	*Logger: 日志对象
func NewLog(out io.Writer, prefix string, flag int) *Logger {
	// 默认 debug打开， calledDepth深度为2, Logger对象调用日志打印方法最多调用两层到达output函数
	l := &Logger{
		out:       out,
		prefix:    prefix,
		flag:      flag,
		file:      nil,
		debug:     true,
		callDepth: 2,
	}
	// 设置log对象 回收资源 析构方法
	runtime.SetFinalizer(l, CleanLog)
	return l
}

// CleanLog 回收日志处理
func CleanLog(log *Logger) {
	log.closeFile()
}

// OutPut 输出日志文件，原方法
func (log *Logger) OutPut(level int, s string) error {
	now := time.Now() // 获取当前时间
	var file string   // 当前调用日志接口的文件名
	var line int      // 当前代码行号

	log.mu.Lock()
	defer log.mu.Unlock()

	// 要显示文件名
	if log.flag&(BitShortFile|BitLongFile) != 0 {
		log.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(log.callDepth)
		if !ok {
			file = "unknown-file"
			line = 0
		}
		log.mu.Lock()
	}

	// 清空缓冲区
	log.buf.Reset()
	// 整理日志头
	log.formatHeader(now, file, line, level)
	// 写日志内容
	log.buf.WriteString(s)
	//补充回车
	if len(s) > 0 && s[len(s)-1] != '\n' {
		log.buf.WriteByte('\n')
	}

	// 将缓冲区内容写到IO输出上
	_, err := log.out.Write(log.buf.Bytes())
	return err
}

func (log *Logger) Debug(v ...interface{}) {
	if !log.debug {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintln(v...))
}

func (log *Logger) Debugf(format string, v ...interface{}) {
	if !log.debug {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintf(format, v...))
}

func (log *Logger) Info(v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintln(v...))
}

func (log *Logger) Infof(format string, v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintf(format, v...))
}

func (log *Logger) Warn(v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintln(v...))
}

func (log *Logger) Warnf(format string, v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintf(format, v...))
}

func (log *Logger) Error(v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintln(v...))
}

func (log *Logger) Errorf(format string, v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintf(format, v...))
}

func (log *Logger) Panic(v ...interface{}) {
	_ = log.OutPut(LogPanic, fmt.Sprintln(v...))
}

func (log *Logger) Panicf(format string, v ...interface{}) {
	_ = log.OutPut(LogPanic, fmt.Sprintf(format, v...))
}

func (log *Logger) Fatal(v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintln(v...))
}

func (log *Logger) Fatalf(format string, v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintf(format, v...))
}

func (log *Logger) Stack(v ...interface{}) {
	s := fmt.Sprintln(v...)
	s += "\n"
	buf := make([]byte, LOG_MAX_BUF)
	n := runtime.Stack(buf, true)
	s += string(buf[:n])
	s += "\n"
	_ = log.OutPut(LogError, s)
}

func (log *Logger) GetLogFlags() int {
	log.mu.Lock()
	defer log.mu.Unlock()

	return log.flag
}

func (log *Logger) SetLogFlags(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.flag = flag
}

func (log *Logger) AddLogFlag(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.flag |= flag
}

func (log *Logger) SetPrefix(prefix string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.prefix = prefix
}

func (log *Logger) SetLogFile(fileDir, fileName string) {
	var file *os.File

	// 创建目录
	err := mkdirLog(fileDir)

	if err != nil {
		fmt.Printf("mkdir log dir error: %v\n", err)
		return
	}

	fullPath := fileDir + "/" + fileName
	if checkFileExist(fullPath) {
		file, err = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0644)
		absPath, _ := filepath.Abs(fullPath)
		fmt.Println("open", absPath)
	} else {
		file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		fmt.Println("create")
	}

	if err != nil {
		fmt.Printf("open log file error: %v\n", err)
		return
	}

	log.mu.Lock()
	defer log.mu.Unlock()

	// 关闭之前绑定的文件
	log.file = file
	log.out = file
}

func (log *Logger) CloseDebug() {
	log.debug = false
}

func (log *Logger) OpenDebug() {
	log.debug = true
}

func (log *Logger) closeFile() {
	if log.file != nil {
		err := log.file.Close()
		if err != nil {
			fmt.Printf("close log file error: %v\n", err)
		}
		log.file = nil
		log.out = os.Stderr
	}
}

// 整理日志头部信息
func (log *Logger) formatHeader(t time.Time, file string, line int, level int) {
	buf := &log.buf

	// 显示头信息
	if log.prefix != "" {
		buf.WriteByte('<')
		buf.WriteString(log.prefix)
		buf.WriteByte('>')
	}

	// 显示时间相关信息
	if log.flag&(BitDate|BitTime|BitMicroSeconds) != 0 {
		// 显示日期
		if log.flag&BitDate != 0 {
			year, month, day := t.Date()
			buf.WriteString(toFixed(year, 4))
			buf.WriteByte('/')
			buf.WriteString(toFixed(int(month), 2))
			buf.WriteByte('/')
			buf.WriteString(toFixed(day, 2))
			buf.WriteByte(' ')
		}

		// 显示时间
		if log.flag&(BitTime|BitMicroSeconds) != 0 {
			hour, min, sec := t.Clock()
			buf.WriteString(toFixed(hour, 2))
			buf.WriteString(":")
			buf.WriteString(toFixed(min, 2))
			buf.WriteString(":")
			buf.WriteString(toFixed(sec, 2))

			// 显示微妙
			if log.flag&BitMicroSeconds != 0 {
				buf.WriteString(".")
				buf.WriteString(toFixed(int(t.UnixMicro()), 6))
			}
			buf.WriteByte(' ')
		}
	}

	// 显示日志级别
	if log.flag&BitLevel != 0 {
		buf.WriteString(levels[level])
	}

	// 显示文件名
	if log.flag&(BitShortFile|BitLongFile) != 0 {
		// 短文件名
		if log.flag&BitShortFile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(toFixed(line, -1))
		buf.WriteString(": ")
	}
}

// 判断日志文件是否存在
func checkFileExist(filename string) bool {
	exists, _ := utils.PathExists(filename)
	return exists
}

// 创建目录
func mkdirLog(dir string) error {
	return utils.Mkdir(dir)
}

// toFixed 将int类型转换为指定长度的字符串
func toFixed(val, length int) string {
	return utils.IntToFixedStr(val, length)
}
