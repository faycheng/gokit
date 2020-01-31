package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrpcPlugin(t *testing.T) {
	plugin := NewGrpcPlugin("../examples/plugin/grpc/echo.bin", "/tmp/gob.echo.socket")
	defer plugin.Close()
	call, err := plugin.Lookup("Echo")
	assert.Nil(t, err)
	_, err = call(context.TODO(), nil)
	assert.Nil(t, err)
}

func BenchmarkGrpcPlugin(b *testing.B) {
	plugin := NewGrpcPlugin("../examples/plugin/grpc/echo.bin", "/tmp/gob.echo.socket")
	defer plugin.Close()
	call, _ := plugin.Lookup("Echo")
	for i := 0; i < b.N; i++ {
		call(context.TODO(), nil)
	}
}
