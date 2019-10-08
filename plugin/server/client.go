package server

import (
	"context"
	"github.com/faycheng/gob/plugin/server/api"
	"google.golang.org/grpc"
)

type PluginClient struct {
	addr   string
	client proto.TaskServiceClient
}

func NewPluginClient(addr string) *PluginClient {
	// TODO: parse the unix socket path
	return &PluginClient{
		addr: addr,
	}
}

func (c *PluginClient) Dial(ctx context.Context) (client proto.TaskServiceClient, err error) {
	// TODO: lock
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	return proto.NewTaskServiceClient(conn), nil
}

func (c *PluginClient) Call(ctx context.Context, method, args string) (err error) {
	if c.client == nil {
		c.client, err = c.Dial(ctx)
		if err != nil {
			return
		}
	}
	req := &proto.CallReq{
		Method: method,
		Args:   args,
	}
	_, err = c.client.Call(ctx, req)
	return
}
