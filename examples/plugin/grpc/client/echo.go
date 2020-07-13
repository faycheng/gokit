package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	plugin2 "github.com/faycheng/gokit/plugin"
	"github.com/sirupsen/logrus"
)

func main() {
	plugin := plugin2.NewGrpcPlugin("../server/echo.bin", "/tmp/gob.echo.socket")
	call, err := plugin.Lookup("Echo")
	if err != nil {
		panic(err)
	}
	reply, err := call(context.TODO(), "hello world")
	logrus.Info("invoke Echo method successfully, reply:%+v err:%+v", reply, err)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	plugin.Close()
	return
}
