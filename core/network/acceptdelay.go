package network

import "time"

const (
	MaxDelay = 1 * time.Second
)

var AcceptDelay = &acceptDelay{}

type acceptDelay struct {
	duration time.Duration
}

func (ad *acceptDelay) GetDuration() time.Duration {
	return ad.duration
}

func (ad *acceptDelay) SetDuration(duration time.Duration) {
	ad.duration = duration
}

func (ad *acceptDelay) Delay() {
	ad.Up()
	ad.doDelay()
}

func (ad *acceptDelay) Up() {
	if ad.duration == 0 {
		ad.duration = 5 * time.Millisecond
		return
	}

	ad.duration *= 2
	if ad.duration > MaxDelay {
		ad.duration = MaxDelay
	}
}

func (ad *acceptDelay) Reset() {
	ad.duration = 0
}

func (ad *acceptDelay) doDelay() {
	if ad.duration <= 0 {
		return
	}

	time.Sleep(ad.duration)
}
