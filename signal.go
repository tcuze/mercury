package main

import (
	"errors"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

type SignalHandler func(...interface{})

var (
	ERR_SIGNAL_NOT_EXIST = errors.New("signal not exist in listener")
	DEFAULT_HANDLER      = func(args interface{}) {}
)

// 非协程安全
type SignalListener struct {
	listener map[os.Signal]SignalHandler
}

func NewSignalListener() *SignalListener {
	return &SignalListener{
		listener: make(map[os.Signal]SignalHandler),
	}
}

func (sl *SignalListener) Register(signal os.Signal, handler SignalHandler) {
	if _, exist := sl.listener[signal]; exist {
		log.Warnf("Signal %s has been registered before, make sure it's in line with expectations", signal.String())
	}
	sl.listener[signal] = handler
}

func (sl *SignalListener) Unregister(signal os.Signal) {
	if _, exist := sl.listener[signal]; exist {
		delete(sl.listener, signal)
	} else {
		log.Warnf("Signal %s has't been registered or removed before", signal.String())
	}
}

func (sl *SignalListener) Handle(signal os.Signal, args ...interface{}) (err error) {
	if _, exist := sl.listener[signal]; exist {
		sl.listener[signal](args)
		return nil
	} else {
		log.Errorf("Signal %s has't been registered or removed before", signal.String())
		return ERR_SIGNAL_NOT_EXIST
	}
}

func (sl *SignalListener) ListenAndServe(channel chan os.Signal, args ...interface{}) {
	go func() {
		for {
			signal.Notify(channel)
			sig := <-channel
			sl.Handle(sig)
		}
	}()
}
