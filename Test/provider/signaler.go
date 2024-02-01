package provider

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	s = make(chan os.Signal, 1)
)

func init() {
	sigs := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGABRT,
	}
	signal.Notify(s, sigs...)
}

// SignalWaitForInterrupt waits until an os.Signal is
// received and returns the result
func SignalWaitForInterrupt() os.Signal {
	return <-s
}
