package task

import (
	"context"
	"strings"
	"time"

	"github.com/faycheng/gokit/signal"
)

type Signal string

const (
	OnStart   Signal = "on_start"
	OnSuccess        = "on_success"
	OnFailure        = "on_failure"
	OnError          = "on_error"
)

type Task interface {
	Name() string
	Connect(Signal, signal.Receiver) error
	Call(c context.Context, args interface{}) error
}

type task struct {
	key       string
	handle    func(c context.Context, args string) error
	onStart   *signal.Signal
	onSuccess *signal.Signal
	onFailure *signal.Signal
	onError   *signal.Signal
}

func NewTask(key string, handle func(c context.Context, args string) error) Task {
	return &task{
		key:       key,
		handle:    handle,
		onStart:   signal.NewSignal(),
		onSuccess: signal.NewSignal(),
		onFailure: signal.NewSignal(),
		onError:   signal.NewSignal(),
	}
}

func (t *task) Name() string {
	return t.key
}

func (t *task) Connect(signal Signal, receiver signal.Receiver) error {
	switch signal {
	case OnStart:
		t.onStart.Connect(receiver)
	case OnSuccess:
		t.onSuccess.Connect(receiver)
	case OnFailure:
		t.onFailure.Connect(receiver)
	case OnError:
		t.onError.Connect(receiver)
	default:
		return nil
	}
	return nil
}

func (t *task) Call(c context.Context, args interface{}) error {
	t.onStart.Send(c, "", args)
	stime := time.Now()
	err := t.handle(c, args.(string))
	if err == nil {
		t.onSuccess.Send(c, "", time.Since(stime))
		return nil
	}
	if strings.HasPrefix(err.Error(), "failure") {
		t.onFailure.Send(c, "", err, args)
		return nil
	}
	t.onError.Send(c, "", err, args)
	return err
}
