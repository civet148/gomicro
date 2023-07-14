package main

import (
	"context"
	"fmt"
	"github.com/civet148/gomicro/v2"
	"github.com/civet148/gomicro/v2/example/echopb"
	"github.com/civet148/log"
	"os"
	"strconv"
)

const (
	SERVICE_NAME         = "echo"
	END_POINTS_MDNS      = "mdns"
	END_POINTS_HTTP_ETCD = "etcd://192.168.2.9:2379"
	//END_POINTS_HTTP_CONSUL   = "consul://192.168.20.108:8500"
	//END_POINTS_TCP_ZOOKEEPER = "zk://192.168.20.108:2181"
	RPC_ADDR = "0.0.0.0:10891" //RPC service listen address
)

type EchoServerImpl struct {
}

func main() {
	var err error
	var weight int
	log.SetLevel("debug")
	var strRpcAddr = RPC_ADDR
	if len(os.Args) == 2 {
		strRpcAddr = fmt.Sprintf("0.0.0.0:%s", os.Args[1])
	} else if len(os.Args) == 3 {
		strRpcAddr = fmt.Sprintf("0.0.0.0:%s", os.Args[1])
		weight, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		log.Debugf("rpc %s weight %v", strRpcAddr, weight)
	}
	ch := make(chan string, 0)
	srv := gomicro.NewServer(END_POINTS_HTTP_ETCD, &gomicro.ServerOption{
		ServiceName: SERVICE_NAME,
		RpcAddr:     strRpcAddr,
		Interval:    3,
		TTL:         10,
		Weight:      weight,
		Metadata: map[string]string{
			"register_name": "echo-server",
		},
	})
	if err := echopb.RegisterEchoServerHandler(srv, new(EchoServerImpl)); err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("micro service starting...")
	if err := srv.Start(); err != nil {
		log.Error(err.Error())
		return
	}
	<-ch
}

func (s *EchoServerImpl) Call(ctx context.Context, ping *echopb.Ping, pong *echopb.Pong) (err error) {
	md := gomicro.FromContext(ctx)
	UserId, _ := md.Get("user_id")
	UserName, _ := md.Get("user_name")
	log.Debugf("md [%+v] req [%+v] user id=[%s] user name [%s]", md, ping, UserId, UserName)
	pong.Text = "Pong"
	return
}
