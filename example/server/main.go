package main

import (
	"context"
	"github.com/civet148/gomicro/v2"
	"github.com/civet148/gomicro/v2/example/echopb"
	"github.com/civet148/log"
	"strings"
	"time"
)

const (
	SERVICE_NAME             = "echo"
	END_POINTS_HTTP_ETCD     = "192.168.20.108:2379"
	END_POINTS_HTTP_CONSUL   = "192.168.20.108:8500"
	END_POINTS_TCP_ZOOKEEPER = "192.168.20.108:2181"
	RPC_ADDR                 = "0.0.0.0:10891" //RPC service listen address
)

type EchoServerImpl struct {
}

func main() {
	log.SetLevel("debug")
	ch := make(chan string, 0)
	srv := NewGoMicroServer(gomicro.RegistryType_ETCD)
	if err := echopb.RegisterEchoServerHandler(srv, new(EchoServerImpl)); err != nil {
		log.Error(err.Error())
		return
	}
	//go-micro v1.16 call srv.Run() v1.18+ call srv.Start()
	if err := srv.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("micro server will deregister after few seconds")
	time.Sleep(15*time.Second)
	if err := srv.Close(); err != nil {
		log.Errorf(err.Error())
	} else {
		log.Infof("micro server deregister ok")
	}
	time.Sleep(10*time.Second)
	srv2 := NewGoMicroServer(gomicro.RegistryType_ETCD)
	if err := echopb.RegisterEchoServerHandler(srv2, new(EchoServerImpl)); err != nil {
		log.Error(err.Error())
		return
	}
	//go-micro v1.16 call srv.Run() v1.18+ call srv.Start()
	if err := srv2.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	<- ch
}

func NewGoMicroServer(typ gomicro.RegistryType) (s *gomicro.GoRPCServer) {
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
	log.Debugf("md [%+v] req [%+v] user id=[%s] user name [%s]", md, ping, UserId, UserName)
	pong.Text = "Pong"
	return
}
