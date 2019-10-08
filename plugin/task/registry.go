package task

import "sync"

type Registry interface {
	Register(string, Task) error
	Get(string) (Task, bool)
}
type taskRegistry struct {
	taskSet sync.Map
}

func NewTaskRegistry() Registry {
	return &taskRegistry{}
}

func (r *taskRegistry) Register(key string, task Task) error {
	r.taskSet.Store(key, task)
	return nil
}

func (r *taskRegistry) Get(key string) (task Task, ok bool) {
	taskIface, ok := r.taskSet.Load(key)
	if !ok {
		return nil, false
	}
	task, ok = taskIface.(Task)
	return
}
