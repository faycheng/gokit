package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin_Run(t *testing.T) {
	config := &PluginConfig{
		Name:  "echo",
		Path:  "../examples/plugin/echo",
		Tasks: []string{"echo"},
	}
	plugin := NewPlugin(config)
	err := plugin.Run()
	assert.Empty(t, err)
	tasks, err := plugin.Tasks()
	assert.Empty(t, err)
	err = tasks["echo"].Call(context.TODO(), `{"msg": "hello world"}`)
	assert.Empty(t, err)
}

// TODO: optimize the performance of task.Run
func BenchmarkPlugin_TaskRun(b *testing.B) {
	config := &PluginConfig{
		Name:  "echo",
		Path:  "../examples/plugin/echo",
		Tasks: []string{"echo"},
	}
	plugin := NewPlugin(config)
	err := plugin.Run()
	if err != nil {
		panic(err)
	}
	tasks, err := plugin.Tasks()
	if err != nil {
		panic(err)
	}
	echoTask := tasks["echo"]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = echoTask.Call(context.TODO(), `{"msg": "hello world"}`)
		if err != nil {
			panic(err)
		}
	}

}
