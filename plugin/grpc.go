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

func (s *grpcServer) Ping(ctx context.Context, req *proto.PingReq) (resp *proto.PingReply, err error) {
	return &proto.PingReply{}, nil
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
	return fmt.Sprintf("type:grpcPlugin pid:%d entrypoint:%s addr:%s", g.cmd.Process.Pid, g.entrypoint, g.client.addr)
}

func (g *grpcPlugin) Close() error {
	err := g.cmd.Process.Kill()
	if err != nil {
		return errors.Wrapf(err, "failed to close grpc plugin server, %s", g)
	}
	logrus.Infof("shutdown grpc plugin successfully, pid:%d entrypoint:%s addr:%s", g.cmd.Process.Pid, g.entrypoint, g.client.addr)
	return nil
}

func NewGrpcPlugin(entrypoint, addr string) Plugin {
	entrypoint, err := filepath.Abs(entrypoint)
	if err != nil {
		panic(errors.Wrapf(err, "failed to expand executable file, addr:%s entrypoint:%s", addr, entrypoint))
	}
	if _, err := os.Stat(entrypoint); err != nil {
		panic(errors.Wrapf(err, "failed to look executable file, addr:%s entrypoint:%s", addr, entrypoint))
	}
	plugin := &grpcPlugin{
		entrypoint: entrypoint,
		client: &grpcClient{
			addr: addr,
		},
	}
	go plugin.run()
	plugin.wait()
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

func (g *grpcPlugin) wait() {
	start := time.Now()
	for time.Since(start) < time.Second*60 {
		conn, err := grpc.Dial(fmt.Sprintf("unix:%s", g.client.addr), grpc.WithInsecure())
		if err != nil {
			logrus.Info(errors.Wrapf(err, "failed to dial grpc plugin server, addr:%s entrypoint:%s", g.client.addr, g.entrypoint))
			time.Sleep(time.Second)
			continue
		}
		client := proto.NewPluginClient(conn)
		_, err = client.Ping(context.TODO(), &proto.PingReq{})
		if err != nil {
			logrus.Info(errors.Wrapf(err, "failed to dial grpc plugin server, addr:%s entrypoint:%s", g.client.addr, g.entrypoint))
			time.Sleep(time.Second)
			continue
		}
		g.client.client = client
		break
	}
	if g.client.client == nil {
		panic(fmt.Sprintf("failed to dial grpc plugin server, addr:%s entrypoint:%s", g.client.addr, g.entrypoint))
	}
	logrus.Infof("start grpc plugin server successfully, addr:%s entrypoint:%s duration:%s", g.client.addr, g.entrypoint, time.Since(start))
}
