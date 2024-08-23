package driven

import (
	"github.com/victorzhou123/simplemq/event"
)

type Publisher interface {
	Publish(event event.Event)
}

type publishImpl struct {
	mq MQ
}

func NewPublish(mq MQ) publishImpl {
	return publishImpl{
		mq: mq,
	}
}

func (p *publishImpl) Publish(event event.Event) {

	if event.Message() == nil {
		return
	}

	p.mq.Push(event.Message())
}
