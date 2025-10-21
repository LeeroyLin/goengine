package iface

type IFlags interface {
	SetString(name, value, usage string)
	SetInt(name string, value int, usage string)
	SetUInt32(name string, value uint32, usage string)
	SetBool(name string, value bool, usage string)

	Parse(args []string) error
	Parsed() bool

	GetString(name, defaultVal string) (string, bool)
	GetBool(name string, defaultVal bool) (bool, bool)
	GetInt(name string, defaultVal int) (int, bool)
	GetUInt32(name string, defaultVal uint32) (uint32, bool)
}
