package main

import (
	"context"

	pluginServer "github.com/faycheng/gokit/plugin/server"
)

var count = 0

func Echo(c context.Context, args interface{}) error {
	count++
	//fmt.Printf("task(echo) is invoking, count(%d)\n", count)
	return nil
}

func main() {
	server := pluginServer.NewPluginServer("test.echo")
	server.Register("echo", Echo)
	err := server.Serve()
	if err != nil {
		panic(err)
	}
}
