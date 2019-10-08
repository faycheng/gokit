package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/faycheng/gob/plugin/server/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
)

type PluginServer struct {
	sync.RWMutex
	addr       string
	grpcServer *grpc.Server
	Tasks      map[string]func(c context.Context, args interface{}) error
}

func (s *PluginServer) Register(task string, handle func(c context.Context, args interface{}) error) {
	s.Tasks[task] = handle
}

func (s *PluginServer) Call(c context.Context, req *proto.CallReq) (resp *proto.CallReply, err error) {
	handle, ok := s.Tasks[req.Method]
	if !ok {
		return nil, fmt.Errorf("method(%s) not found", req.Method)
	}
	var args interface{}
	err = json.Unmarshal([]byte(req.Args), &args)
	if err != nil {
		logrus.Errorf("bad Request(%+v)", req)
		return
	}
	return &proto.CallReply{}, handle(c, args)
}

func NewPluginServer(name string) *PluginServer {
	return &PluginServer{
		// TODO: /var/run
		addr:  fmt.Sprintf("/tmp/gob/%s.socket", name),
		Tasks: make(map[string]func(c context.Context, args interface{}) error),
	}
}

func (s *PluginServer) Serve() (err error) {
	if _, err := os.Stat(s.addr); err == nil {
		logrus.Warnf("unix socket(%s) with same name already exists", s.addr)
		err = os.Remove(s.addr)
		if err != nil {
			return err
		}
	}
	listener, err := net.Listen("unix", s.addr)
	if err != nil {
		return
	}
	s.grpcServer = grpc.NewServer()
	proto.RegisterTaskServiceServer(s.grpcServer, s)
	return s.grpcServer.Serve(listener)
}
