package signal

import (
	"context"
	"sync"
)

const (
	ALL = "ALL"
)

type Receiver func(c context.Context, args ...interface{}) error

type Signal struct {
	Receivers  sync.Map
	StrictMode bool
}

func NewSignal(strictMode bool) *Signal {
	return &Signal{
		StrictMode: strictMode,
	}
}

func (s *Signal) Connect(receiver Receiver, sender string) error {
	var receivers []Receiver
	if sender == "" {
		sender = ALL
	}
	receiversIface, ok := s.Receivers.Load(sender)
	if !ok {
		receivers = make([]Receiver, 0)
	} else {
		receivers = receiversIface.([]Receiver)
	}
	receivers = append(receivers, receiver)
	s.Receivers.Store(sender, receivers)
	return nil
}

func (s *Signal) SendTo(c context.Context, sender string, args ...interface{}) error {
	receiversIface, ok := s.Receivers.Load(sender)
	if !ok {
		return nil
	}
	receivers := receiversIface.([]Receiver)
	for _, receiver := range receivers {
		err := receiver(c, args...)
		if err != nil && s.StrictMode {
			return err
		}
	}
	return nil
}

func (s *Signal) Send(c context.Context, sender string, args ...interface{}) error {
	var err error
	err = s.SendTo(c, ALL, args...)
	if err != nil && s.StrictMode {
		return err
	}
	if sender == "" || sender == ALL {
		return nil
	}
	err = s.SendTo(c, sender, args...)
	return err
}
