package main

import (
	"context"
	"fmt"
	"time"

	"github.com/faycheng/gokit/plugin"
)

func Echo(ctx context.Context, req interface{}) (reply interface{}, err error) {
	fmt.Println("hello world")
	return nil, nil
}

func main() {
	server := plugin.NewGrpcPluginServer("/tmp/gob.echo.socket")
	server.Register("Echo", Echo)
	time.Sleep(time.Second * 3)
	err := server.Serve()
	if err != nil {
		panic(err)
	}
}
