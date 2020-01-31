package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/faycheng/gokit/plugin/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcClient struct {
	addr   string
	client proto.PluginClient
}

func (c *grpcClient) Get(ctx context.Context, name string) (reply interface{}, err error) {
	req := &proto.GetReq{
		Name: name,
	}
	_, err = c.client.Get(ctx, req)
	return nil, err
}

func (c *grpcClient) Call(ctx context.Context, name string, req interface{}) (reply interface{}, err error) {
	args := []byte("{}")
	if req != nil {
		args, _ = json.Marshal(req)
	}
	callReq := &proto.CallReq{
		Name: name,
		Args: args,
	}
	_, err = c.client.Call(ctx, callReq)
	return
}

type grpcServer struct {
	sync.RWMutex
	addr       string
	grpcServer *grpc.Server
	calls      map[string]Call
}

func (s *grpcServer) Register(name string, call Call) {
	s.calls[name] = call
}

func (s *grpcServer) Call(ctx context.Context, req *proto.CallReq) (resp *proto.CallReply, err error) {
	call, ok := s.calls[req.Name]
	if !ok {
		return nil, fmt.Errorf("grpc call handler not found, name:%s", req.Name)
	}
	var args interface{}
	err = json.Unmarshal([]byte(req.Args), &args)
	if err != nil {
		logrus.Errorf("bad Request(%+v)", req)
		return
	}
	_, err = call(ctx, args)
	return &proto.CallReply{}, err
}

func (s *grpcServer) Get(ctx context.Context, req *proto.GetReq) (resp *proto.GetReply, err error) {
	_, ok := s.calls[req.Name]
	if !ok {
		return nil, fmt.Errorf("grpc call handler not found, name:%s", req.Name)
	}
	return &proto.GetReply{}, nil
}

func NewGrpcPluginServer(addr string) *grpcServer {
	return &grpcServer{
		// TODO: /var/run
		addr:  addr,
		calls: make(map[string]Call),
	}
}

func (s *grpcServer) Serve() (err error) {
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
	proto.RegisterPluginServer(s.grpcServer, s)
	return s.grpcServer.Serve(listener)
}

type grpcPlugin struct {
	entrypoint string
	cmd        *exec.Cmd
	client     *grpcClient
}

func (g *grpcPlugin) Lookup(name string) (call Call, err error) {
	_, err = g.client.Get(context.TODO(), name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find grpc method, name:%s", name)
	}
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		return g.client.Call(ctx, name, req)
	}, nil
}

func (g *grpcPlugin) String() string {
	return fmt.Sprintf("type:soPlugin pid:%d entrypoint:%s addr:%s", g.cmd.Process.Pid, g.entrypoint, g.client.addr)
}

func (g *grpcPlugin) Close() error {
	err := g.cmd.Process.Kill()
	if err != nil {
		return errors.Wrapf(err, "failed to close grpc plugin server, %s", g)
	}
	return nil
}

func NewGrpcPlugin(entrypoint, addr string) Plugin {
	//var err error
	//entrypoint, err = filepath.Abs(entrypoint)
	entrypoint, err := filepath.Abs(entrypoint)
	if err != nil {
		panic(errors.Wrapf(err, "failed to look executable file, addr:%s entrypoint:%s", addr, entrypoint))
	}
	plugin := &grpcPlugin{
		entrypoint: entrypoint,
	}
	go plugin.run()
	time.Sleep(100 * time.Millisecond)
	conn, err := grpc.Dial(fmt.Sprintf("unix:%s", addr), grpc.WithInsecure())
	if err != nil {
		panic(errors.Wrapf(err, "failed to dial grpc plugin server, addr:%s entrypoint:%s", addr, entrypoint))
	}
	plugin.client = &grpcClient{
		addr:   addr,
		client: proto.NewPluginClient(conn),
	}
	return plugin
}

func (g *grpcPlugin) run() {
	logrus.Infof("start grpc plugin server, entrypoint:%s", g.entrypoint)
	cmd := exec.Command(g.entrypoint)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	g.cmd = cmd
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("exit grpc plugin server, entrypoint:%s err:%+v", g.entrypoint, err)
		return
	}
	logrus.Infof("exit grpc plugin server, entrypoint:%s", g.entrypoint)
}
