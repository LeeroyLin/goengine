package cli

func (c *Cmd) AddStringFlag(name, shorthand string, value string, usage string) *Cmd {
	v := value

	if c.strMap == nil {
		c.strMap = make(map[*string]string)
	}

	c.strMap[&v] = value

	c.root.Flags().StringVarP(&v, name, shorthand, value, usage)

	return c
}

func (c *Cmd) AddBoolFlag(name, shorthand string, value bool, usage string) *Cmd {
	v := value

	if c.boolMap == nil {
		c.boolMap = make(map[*bool]bool)
	}

	c.boolMap[&v] = value

	c.root.Flags().BoolVarP(&v, name, shorthand, value, usage)

	return c
}

func (c *Cmd) AddIntFlag(name, shorthand string, value int, usage string) *Cmd {
	v := value

	if c.intMap == nil {
		c.intMap = make(map[*int]int)
	}

	c.intMap[&v] = value

	c.root.Flags().IntVarP(&v, name, shorthand, value, usage)

	return c
}

func (c *Cmd) AddInt64Flag(name, shorthand string, value int64, usage string) *Cmd {
	v := value

	if c.int64Map == nil {
		c.int64Map = make(map[*int64]int64)
	}

	c.int64Map[&v] = value

	c.root.Flags().Int64VarP(&v, name, shorthand, value, usage)

	return c
}

func (c *Cmd) AddUint32Flag(name, shorthand string, value uint32, usage string) *Cmd {
	v := value

	if c.uint32Map == nil {
		c.uint32Map = make(map[*uint32]uint32)
	}

	c.uint32Map[&v] = value

	c.root.Flags().Uint32VarP(&v, name, shorthand, value, usage)

	return c
}

func (c *Cmd) resetAllFlags() {
	if c.strMap != nil {
		for p, defaultV := range c.strMap {
			*p = defaultV
		}
	}

	if c.boolMap != nil {
		for p, defaultV := range c.boolMap {
			*p = defaultV
		}
	}

	if c.intMap != nil {
		for p, defaultV := range c.intMap {
			*p = defaultV
		}
	}

	if c.int64Map != nil {
		for p, defaultV := range c.int64Map {
			*p = defaultV
		}
	}

	if c.uint32Map != nil {
		for p, defaultV := range c.uint32Map {
			*p = defaultV
		}
	}
}
