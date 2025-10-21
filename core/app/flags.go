package app

import (
	"flag"
)

type Flags struct {
	flagSet *flag.FlagSet
}

func NewFlags() *Flags {
	af := &Flags{
		flagSet: &flag.FlagSet{},
	}

	return af
}

func (flags *Flags) SetString(name, value, usage string) {
	flags.flagSet.String(name, value, usage)
}

func (flags *Flags) SetInt(name string, value int, usage string) {
	flags.flagSet.Int(name, value, usage)
}

func (flags *Flags) SetUInt32(name string, value uint32, usage string) {
	flags.flagSet.Uint(name, uint(value), usage)
}

func (flags *Flags) SetBool(name string, value bool, usage string) {
	flags.flagSet.Bool(name, value, usage)
}

func (flags *Flags) Parse(args []string) error {
	return flags.flagSet.Parse(args)
}

func (flags *Flags) Parsed() bool {
	return flags.Parsed()
}

func (flags *Flags) GetString(name, defaultVal string) (string, bool) {
	f := flags.flagSet.Lookup(name)
	if f == nil || f.Value == nil {
		return defaultVal, false
	}

	return f.Value.String(), true
}

func (flags *Flags) GetBool(name string, defaultVal bool) (bool, bool) {
	f := flags.flagSet.Lookup(name)
	if f == nil || f.Value == nil {
		return defaultVal, false
	}

	return f.Value.(flag.Getter).Get().(bool), true
}

func (flags *Flags) GetInt(name string, defaultVal int) (int, bool) {
	f := flags.flagSet.Lookup(name)
	if f == nil || f.Value == nil {
		return defaultVal, false
	}

	return f.Value.(flag.Getter).Get().(int), true
}

func (flags *Flags) GetUInt32(name string, defaultVal uint32) (uint32, bool) {
	f := flags.flagSet.Lookup(name)
	if f == nil || f.Value == nil {
		return defaultVal, false
	}

	return uint32(f.Value.(flag.Getter).Get().(uint)), true
}
