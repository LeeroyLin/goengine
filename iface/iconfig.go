package iface

type IConfig interface {
	Setup(confFilePath string, flags IFlags)
	LoadFromFile(confFilePath string)
	InitFlags(flags IFlags)
	ParseFlags(flags IFlags)
	GetLogStr() string
}
