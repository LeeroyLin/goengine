package closer

import (
	"os"
	"os/signal"
	"syscall"
)

type SigCloser struct {
	sigChan   chan os.Signal
	CloseChan chan interface{}
}

func NewSigCloser() *SigCloser {
	return &SigCloser{
		sigChan:   make(chan os.Signal, 1),
		CloseChan: make(chan interface{}),
	}
}

func (c *SigCloser) Listen(final func()) {
	// 监听 ctrl+c 终止信号 挂断信号
	signal.Notify(c.sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case <-c.sigChan:
			final()
			return
		case <-c.CloseChan:
			return
		}
	}
}

func (c *SigCloser) Close() {
	select {
	case <-c.CloseChan:
		return
	default:
		close(c.CloseChan)
	}
}
