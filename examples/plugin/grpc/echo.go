package main

import (
	"context"
	"fmt"

	"github.com/faycheng/gokit/plugin"
)

func Echo(ctx context.Context, req interface{}) (reply interface{}, err error) {
	fmt.Println("hello world")
	return nil, nil
}

func main() {
	server := plugin.NewGrpcPluginServer("/tmp/gob.echo.socket")
	server.Register("Echo", Echo)
	err := server.Serve()
	if err != nil {
		panic(err)
	}
}
