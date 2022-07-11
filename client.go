package gomicro

import (
	"context"
	"github.com/micro/go-micro/v2/client"
)

func (c *GoRPCClient) Init(opts ...client.Option) error {
	return c.client.Init(opts...)
}

func (c *GoRPCClient) Options() client.Options {
	return c.client.Options()
}

func (c *GoRPCClient) NewMessage(topic string, msg interface{}, opts ...client.MessageOption) client.Message {
	return c.client.NewMessage(topic, msg, opts...)
}

func (c *GoRPCClient) NewRequest(service, endpoint string, req interface{}, reqOpts ...client.RequestOption) client.Request {
	return c.client.NewRequest(service, endpoint, req, reqOpts...)
}

func (c *GoRPCClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return c.client.Call(ctx, req, rsp, opts...)
}

func (c *GoRPCClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	return c.client.Stream(ctx, req, opts...)
}

func (c *GoRPCClient) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	return c.client.Publish(ctx, msg, opts...)
}

func (c *GoRPCClient) String() string {
	return c.client.String()
}