package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/faycheng/gokit/plugin/server"
	"github.com/faycheng/gokit/plugin/task"
	"github.com/sirupsen/logrus"
)

type Plugin interface {
	Tasks() (map[string]task.Task, error)
	Run() error
}

type plugin struct {
	name    string
	path    string
	tasks   []string
	running int32
	client  *server.PluginClient
}

type PluginConfig struct {
	Name  string `json:"Name"`
	Path  string
	Tasks []string `json:"Tasks"`
}

func NewPlugin(config *PluginConfig) Plugin {
	return &plugin{
		name:  config.Name,
		path:  config.Path,
		tasks: config.Tasks,
	}
}

// TODO: thread-safe
func (p *plugin) Tasks() (map[string]task.Task, error) {
	if atomic.LoadInt32(&p.running) == 0 {
		return nil, fmt.Errorf("plugin server isn't running, please invoke plugin.Run before fetching tasks")
	}
	client := server.NewPluginClient("unix:/tmp/gob/test.echo.socket")
	tasks := make(map[string]task.Task)
	for _, name := range p.tasks {
		handle := func(c context.Context, args string) error {
			// TODO: register on_start, on_success... events for collecting metrics
			return client.Call(c, name, args)
		}
		t := task.NewTask(name, handle)
		//counter, _ := factory.NewCounter()
		//gauge, _ := factory.NewGauge()
		//t.Connect(task.OnSuccess, func(c context.Context, args ...interface{}) error {
		//	logrus.Infof("task.OnStart(%+v)", args)
		//	counter.Add(1)
		//	gauge.Add(int64(args[0].(time.Duration) / time.Millisecond))
		//	//spew.Dump(counter)
		//	//spew.Dump(gauge)
		//	return nil
		//})
		tasks[name] = t
	}
	return tasks, nil
}

func (p *plugin) Run() error {
	entrypoint := fmt.Sprintf("%s/entrypoint", p.path)
	_, err := exec.LookPath(entrypoint)
	if err != nil {
		return err
	}
	cmd := exec.Command(entrypoint)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	atomic.StoreInt32(&p.running, 1)
	go func() {
		err := cmd.Run()
		if err != nil {
			logrus.Errorf("plugin(%s/entrypoint.sh) exit with err: %+v", p.path, err)
			return
		}
		logrus.Infof("plugin(%s,entrypoint.sh) exit", p.path)
	}()
	time.Sleep(1 * time.Second)
	return nil
}
