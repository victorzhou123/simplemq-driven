package main

import (
	"sync"

	"github.com/victorzhou123/simplemq/event"
)

type Consumer interface {
	Consume(event.Event)
	Topics() []string
}

type SubsMap map[string][]Consumer

func newEmptySubsMap() SubsMap {
	return make(map[string][]Consumer, 0)
}

func (s SubsMap) Add(sub Consumer) {

	if sub == nil || len(sub.Topics()) == 0 {
		return
	}

	for _, topic := range sub.Topics() {

		_, ok := s[topic]
		if !ok {
			s[topic] = []Consumer{sub}

			continue
		}

		s[topic] = append(s[topic], sub)
	}
}

func (s SubsMap) Find(topic string) ([]Consumer, bool) {

	arr, ok := s[topic]
	if !ok {
		return nil, false
	}

	if len(arr) == 0 {
		return nil, false
	}

	return arr, true
}

type Distributer interface {
	Distribute(event.Event)
	Subscribe(Consumer)
}

type distributerImpl struct {
	m    sync.Mutex
	subs SubsMap
}

func NewDistributerImpl() Distributer {
	return &distributerImpl{
		subs: newEmptySubsMap(),
	}
}

func (s *distributerImpl) Subscribe(sub Consumer) {

	s.m.Lock()
	defer s.m.Unlock()

	s.subs.Add(sub)
}

func (s *distributerImpl) Distribute(e event.Event) {

	subs, ok := s.subs.Find(e.Topic())
	if !ok {
		return
	}

	for _, sub := range subs {
		go sub.Consume(e)
	}
}
