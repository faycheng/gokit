package signal

import (
	"context"
	"sync"
)

type Receiver func(c context.Context, args ...interface{}) error

type Signal struct {
	sync.RWMutex
	Receivers  []Receiver
	StrictMode bool
}

func NewSignal() *Signal {
	return &Signal{
		Receivers: make([]Receiver, 0),
	}
}

func (s *Signal) Connect(receiver Receiver) error {
	s.Receivers = append(s.Receivers, receiver)
	return nil
}

func (s *Signal) Send(c context.Context, args ...interface{}) error {
	s.Lock()
	defer s.Unlock()
	for _, receiver := range s.Receivers {
		err := receiver(c, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
