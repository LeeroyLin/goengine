package iface

type IConf interface {
	LoadFromFile(confFilePath string)
	InitFlags(flags IFlags)
	ParseFlags(flags IFlags)
	GetLogStr() string
}
