package main

import (
	"time"

	"github.com/victorzhou123/simplemq/event"
)

type MQ interface {
	Push(*event.Message)
	Pop() event.Message
	HasMsg() bool
}

type watcher struct {
	mq        MQ
	subscribe Distributer
	period    time.Duration
}

func NewWatcher(mq MQ, sub Distributer) watcher {
	return watcher{
		mq:        mq,
		subscribe: sub,
		period:    20 * time.Millisecond,
	}
}

func (w *watcher) Watch() {

	ticker := time.NewTicker(w.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if w.mq.HasMsg() {
				msg := w.mq.Pop()
				w.subscribe.Distribute(msg)
			}
		}
	}
}
