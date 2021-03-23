package signal

import (
	"os"
	"os/signal"
	"time"
)

func Register(sig os.Signal, handler signalFunc) {
	signalMaps[sig.String()] = handler
	signal.Notify(sigChan, sig)
}

type signalFunc func(sig os.Signal)

var sigChan chan os.Signal

func run() {
	go func() {
		for true {
			select {
			case sig := <-sigChan:
				method := signalMaps[sig.String()]
				if method != nil {
					method(sig)
				}

			default:
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}()
}

var signalMaps map[string]signalFunc

func init() {
	signalMaps = make(map[string]signalFunc)
	sigChan = make(chan os.Signal, 3)
	run()
}
