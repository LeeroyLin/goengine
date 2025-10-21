package conf

import (
	"bytes"
	"encoding/json"
	"engine/core/elog"
	"engine/core/utils"
	"engine/iface"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Conf struct {
	Name      string // 名字
	Desc      string // 额外描述
	IPVersion string // 主机ip版本：tcp,tcp4,tcp6
	IP        string // 主机ip
	Port      int    // 主机端口号

	LogDir   string // 日志文件目录
	LogFile  string // 日志文件名
	LogDebug bool   // 开启日志调试

	MaxConn int // 最大连接数

	WorkerPoolSize    uint32 // 工作池数量
	MaxWorkerTaskLen  uint32 // worker最大任务容量
	MaxPacketSize     uint32 // 最大包长度
	MaxMsgBuffChanLen uint32 // 最大消息队列通道容量
}

func NewConf() *Conf {
	c := &Conf{
		Name: "UnknownServer",
		Desc: "NoDesc",
		IP:   "0.0.0.0",
		Port: 8999,
	}

	return c
}

// Setup 装载配置
func (c *Conf) Setup(confFilePath string, flags iface.IFlags) {
	// 加载配置文件
	c.LoadFromFile(confFilePath)

	// 初始化命令行参数
	c.InitFlags(flags)

	// 读取命令行参数
	cliArgs := os.Args[1:]
	err := flags.Parse(cliArgs)
	if err != nil {
		panic(err)
		return
	}

	// 处理命令行参数
	c.ParseFlags(flags)

	// 日志文件
	if c.LogFile != "" {
		elog.SetLogFile(c.LogDir, c.LogFile)
	}

	// 日志调试
	if c.LogDebug {
		elog.OpenDebug()
	} else {
		elog.CloseDebug()
	}
}

// LoadFromFile 从文件加载配置
func (c *Conf) LoadFromFile(confFilePath string) {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	if utils.IsEmptyOrWhitespace(confFilePath) {
		confFilePath = pwd + "/conf/config.json"
	}

	// 检测路径
	ok, err := utils.PathExists(confFilePath)

	// 路径不存在
	if !ok {
		panic(err)
	}

	// 读取配置
	data, err := os.ReadFile(confFilePath)
	if err != nil {
		panic(err)
	}

	// 将json数据解析到struct中
	err = json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
}

func (c *Conf) InitFlags(flags iface.IFlags) {
	// 获取反射值对象
	val := reflect.ValueOf(c)

	// 如果是指针类型，获取其指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取结构体类型
	typ := val.Type()

	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typeField := typ.Field(i)
		fieldName := typeField.Name
		lowerName := getCmdName(fieldName)
		fieldValue := valField.Interface()

		kind := typeField.Type.Kind()

		if kind == reflect.String {
			flags.SetString(lowerName, fieldValue.(string), fieldName)
		} else if kind == reflect.Bool {
			flags.SetBool(lowerName, fieldValue.(bool), fieldName)
		} else if kind == reflect.Int {
			flags.SetInt(lowerName, fieldValue.(int), fieldName)
		} else if kind == reflect.Uint32 {
			flags.SetUInt32(lowerName, fieldValue.(uint32), fieldName)
		}
	}

	flags.SetString("test.v", "", "test")
	flags.SetBool("test.paniconexit0", false, "test")
	flags.SetString("test.run", "", "test")
	flags.SetString("test.timeout", "", "test")
}

func (c *Conf) ParseFlags(flags iface.IFlags) {
	// 获取反射值对象
	val := reflect.ValueOf(c)

	// 如果是指针类型，获取其指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取结构体类型
	typ := val.Type()

	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)

		if !valField.CanSet() {
			continue
		}

		typeField := typ.Field(i)
		fieldName := typeField.Name
		lowerName := getCmdName(fieldName)
		fieldValue := valField.Interface()

		kind := typeField.Type.Kind()

		if kind == reflect.String {
			v, _ := flags.GetString(lowerName, fieldValue.(string))
			valField.SetString(v)
		} else if kind == reflect.Bool {
			v, _ := flags.GetBool(lowerName, fieldValue.(bool))
			valField.SetBool(v)
		} else if kind == reflect.Int {
			v, _ := flags.GetInt(lowerName, fieldValue.(int))
			valField.SetInt(int64(v))
		} else if kind == reflect.Uint32 {
			v, _ := flags.GetUInt32(lowerName, fieldValue.(uint32))
			valField.SetUint(uint64(v))
		}
	}
}

func (c *Conf) GetLogStr() string {
	var buf bytes.Buffer

	// 获取反射值对象
	val := reflect.ValueOf(c)

	// 如果是指针类型，获取其指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取结构体类型
	typ := val.Type()

	buf.WriteString("\n")

	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := field.Interface()

		buf.WriteString(fmt.Sprintf("%s:%v,\n", fieldName, fieldValue))
	}

	return buf.String()
}

func getCmdName(fieldName string) string {
	res := strings.ReplaceAll(fieldName, "_", "")
	return strings.ToLower(res)
}
