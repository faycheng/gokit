package signal

import "sync"

type NameSpace struct {
	Signals    sync.Map
	StrictMode bool
}

func NewNamespace(strictMode bool) *NameSpace {
	return &NameSpace{
		StrictMode: strictMode,
	}
}

func (n *NameSpace) Signal(name string) *Signal {
	var signal *Signal
	signalIface, ok := n.Signals.Load(name)
	if !ok {
		signal = NewSignal(n.StrictMode)
	} else {
		signal = signalIface.(*Signal)
	}
	return signal
}
