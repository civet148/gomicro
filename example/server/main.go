package main

import (
	"context"
	"github.com/civet148/gomicro"
	"github.com/civet148/gomicro/example/echopb"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/server"
	"strings"
)

const (
	SERVICE_NAME             = "echo"
	END_POINTS_HTTP_ETCD     = "http://127.0.0.1:2379"
	END_POINTS_HTTP_CONSUL   = "http://127.0.0.1:8500"
	END_POINTS_TCP_ZOOKEEPER = "127.0.0.1:2181"
	RPC_ADDR                 = "127.0.0.1:8899" //RPC service listen address
)

type EchoServerImpl struct {
}

func main() {
	ch := make(chan bool, 1)
	srv := NewGoMicroServer(gomicro.RegistryType_MDNS)
	if err := echopb.RegisterEchoServerHandler(srv, new(EchoServerImpl)); err != nil {
		log.Error(err.Error())
		return
	}
	//go-micro v1.16 call srv.Run() v1.18 call srv.Start()
	if err := srv.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	<-ch //block infinite
}

func NewGoMicroServer(typ gomicro.RegistryType) (s server.Server) {
	var g *gomicro.GoRPC
	var endPoints []string

	g = gomicro.NewGoRPC(typ)
	switch typ {
	case gomicro.RegistryType_MDNS:
	case gomicro.RegistryType_ETCD:
		endPoints = strings.Split(END_POINTS_HTTP_ETCD, ",")
	//case gomicro.RegistryType_CONSUL:
	//	endPoints = strings.Split(END_POINTS_HTTP_CONSUL, ",")
	//case gomicro.RegistryType_ZOOKEEPER:
	//	endPoints = strings.Split(END_POINTS_TCP_ZOOKEEPER, ",")
	}

	return g.NewServer(&gomicro.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   endPoints,
	})
}

func (s *EchoServerImpl) Call(ctx context.Context, ping *echopb.Ping, pong *echopb.Pong) (err error) {
	md := gomicro.FromContext(ctx)
	UserId, _ := md.Get("user_id")
	UserName, _ := md.Get("user_name")
	log.Infof("md [%+v] req [%+v] user id=[%s] user name [%s]", md, ping, UserId, UserName)
	pong.Text = "Pong"
	return
}
