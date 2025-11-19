package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/flags"
	"github.com/LeeroyLin/goengine/core/utils"
	"os"
	"reflect"
	"strings"
)

var buf bytes.Buffer

type ConfBase struct {
	Flags *flags.Flags
	ConfBasicPattern
	ConfLogPattern
}

func NewConfBase() ConfBase {
	return ConfBase{
		Flags:            flags.NewFlags(),
		ConfBasicPattern: ConfBasicPattern{},
		ConfLogPattern:   ConfLogPattern{},
	}
}

type ConfBasicPattern struct {
	Name string // 名字
	Desc string // 额外描述
}

type ConfLogPattern struct {
	LogDir   string // 日志文件目录
	LogFile  string // 日志文件名
	LogDebug bool   // 开启日志调试
}

type ConfHttpServicePattern struct {
	HttpIP       string // 主机ip
	HttpPort     int    // 主机端口号
	Https        bool
	HttpCertFile string
	HttpKeyFile  string
}

func NewConfHttpServicePattern() ConfHttpServicePattern {
	return ConfHttpServicePattern{}
}

type ConfNetServicePattern struct {
	IPVersion string // 主机ip版本：tcp,tcp4,tcp6
	IP        string // 主机ip
	Port      int    // 主机端口号

	MaxConn           int    // 最大连接数
	WorkerPoolSize    uint32 // 工作池数量
	MaxWorkerTaskLen  uint32 // worker最大任务容量
	MaxPacketSize     uint32 // 最大包长度
	MaxMsgBuffChanLen uint32 // 最大消息队列通道容量
}

func NewConfNetServicePattern() ConfNetServicePattern {
	return ConfNetServicePattern{}
}

type ConfETCDPattern struct {
	ETCDServerId  string
	ETCDEndpoints []string
	ETCDTTL       int64
}

func NewConfETCDPattern() ConfETCDPattern {
	return ConfETCDPattern{}
}

// Setup 装载配置
func (c *ConfBase) Setup(child interface{}, confFilePath string) {
	// 加载配置文件
	c.LoadFromFile(child, confFilePath)

	// 初始化命令行参数
	c.initFlags(child)

	// 初始化测试命令行参数
	c.initTestFlags()

	// 读取命令行参数
	cliArgs := os.Args[1:]
	err := c.Flags.Parse(cliArgs)
	if err != nil {
		panic(err)
		return
	}

	// 处理命令行参数
	c.parseFlags(child)

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
func (c *ConfBase) LoadFromFile(child interface{}, confFilePath string) {
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
	err = json.Unmarshal(data, child)
	if err != nil {
		panic(err)
	}
}

func (c *ConfBase) initFlags(child interface{}) {
	// 获取反射值对象
	val := reflect.ValueOf(child)

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

		if fieldName == "Flags" {
			continue
		}

		kind := typeField.Type.Kind()

		if kind == reflect.String {
			c.Flags.SetString(lowerName, fieldValue.(string), fieldName)
		} else if kind == reflect.Bool {
			c.Flags.SetBool(lowerName, fieldValue.(bool), fieldName)
		} else if kind == reflect.Int {
			c.Flags.SetInt(lowerName, fieldValue.(int), fieldName)
		} else if kind == reflect.Uint32 {
			c.Flags.SetUInt32(lowerName, fieldValue.(uint32), fieldName)
		} else if kind == reflect.Struct {
			if typeField.Type == reflect.TypeOf(ConfBase{}) {
				if valField.CanSet() {
					ptr := valField.Addr().Interface().(*ConfBase)

					c.initFlags(ptr)
				}
			}
		}
	}
}

func (c *ConfBase) initTestFlags() {
	c.Flags.SetString("test.v", "", "test")
	c.Flags.SetBool("test.paniconexit0", false, "test")
	c.Flags.SetString("test.run", "", "test")
	c.Flags.SetString("test.timeout", "", "test")
}

func (c *ConfBase) parseFlags(child interface{}) {
	// 获取反射值对象
	val := reflect.ValueOf(child)

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

		if fieldName == "Flags" {
			continue
		}

		kind := typeField.Type.Kind()

		if kind == reflect.String {
			v, _ := c.Flags.GetString(lowerName, fieldValue.(string))
			valField.SetString(v)
		} else if kind == reflect.Bool {
			v, _ := c.Flags.GetBool(lowerName, fieldValue.(bool))
			valField.SetBool(v)
		} else if kind == reflect.Int {
			v, _ := c.Flags.GetInt(lowerName, fieldValue.(int))
			valField.SetInt(int64(v))
		} else if kind == reflect.Uint32 {
			v, _ := c.Flags.GetUInt32(lowerName, fieldValue.(uint32))
			valField.SetUint(uint64(v))
		} else if kind == reflect.Struct {
			if typeField.Type == reflect.TypeOf(ConfBase{}) {
				if valField.CanSet() {
					ptr := valField.Addr().Interface().(*ConfBase)

					c.parseFlags(ptr)
				}
			}
		}
	}
}

func (c *ConfBase) GetLogStr(child interface{}) string {
	buf.Reset()

	buf.WriteString("\n[Conf] ====================\n")

	c.recordLogStr(child)

	buf.WriteString("[Conf] ====================\n")

	return buf.String()
}

func (c *ConfBase) recordLogStr(child interface{}) {
	// 获取反射值对象
	val := reflect.ValueOf(child)

	// 如果是指针类型，获取其指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取结构体类型
	typ := val.Type()

	// 遍历结构体的所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldTyp := typ.Field(i)

		if fieldTyp.Type.Kind() == reflect.Struct {
			c.recordLogStr(field.Interface())
			continue
		}

		fieldName := fieldTyp.Name
		fieldValue := field.Interface()

		if fieldName == "Flags" {
			continue
		}

		buf.WriteString("    ")
		buf.WriteString(fmt.Sprintf("%s:%v,\n", fieldName, fieldValue))
	}
}

func getCmdName(fieldName string) string {
	res := strings.ReplaceAll(fieldName, "_", "")
	return strings.ToLower(res)
}
