package utils

import "time"

type PowDelay struct {
	start    time.Duration
	multi    float32
	max      time.Duration
	duration time.Duration
}

func NewPowDelay(start, max time.Duration, multi float32) *PowDelay {
	d := &PowDelay{
		start:    start,
		max:      max,
		multi:    multi,
		duration: start,
	}

	return d
}

func (pd *PowDelay) Delay() {
	pd.Up()
	pd.doDelay()
}

func (pd *PowDelay) Up() {
	d := int64(pd.duration)
	c := float32(d) * pd.multi
	pd.duration = time.Duration(c)
	if pd.duration > pd.max {
		pd.duration = pd.max
	}
}

func (pd *PowDelay) Reset() {
	pd.duration = pd.start
}

func (pd *PowDelay) doDelay() {
	if pd.duration <= 0 {
		return
	}

	time.Sleep(pd.duration)
}
