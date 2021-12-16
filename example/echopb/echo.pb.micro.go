// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: echo.proto

package echopb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for EchoServer service

func NewEchoServerEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for EchoServer service

type EchoServerService interface {
	Call(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error)
}

type echoServerService struct {
	c    client.Client
	name string
}

func NewEchoServerService(name string, c client.Client) EchoServerService {
	return &echoServerService{
		c:    c,
		name: name,
	}
}

func (c *echoServerService) Call(ctx context.Context, in *Ping, opts ...client.CallOption) (*Pong, error) {
	req := c.c.NewRequest(c.name, "EchoServer.Call", in)
	out := new(Pong)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EchoServer service

type EchoServerHandler interface {
	Call(context.Context, *Ping, *Pong) error
}

func RegisterEchoServerHandler(s server.Server, hdlr EchoServerHandler, opts ...server.HandlerOption) error {
	type echoServer interface {
		Call(ctx context.Context, in *Ping, out *Pong) error
	}
	type EchoServer struct {
		echoServer
	}
	h := &echoServerHandler{hdlr}
	return s.Handle(s.NewHandler(&EchoServer{h}, opts...))
}

type echoServerHandler struct {
	EchoServerHandler
}

func (h *echoServerHandler) Call(ctx context.Context, in *Ping, out *Pong) error {
	return h.EchoServerHandler.Call(ctx, in, out)
}
