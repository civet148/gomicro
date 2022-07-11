package main

import (
	"fmt"
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
)

func main() {

	c := NewGoMicroClient(gomicro.RegistryType_ETCD)
	service := echopb.NewEchoServerService(SERVICE_NAME, c.Client)

	for i := 0; i < 200; i++ {
		ctx := gomicro.NewContext(map[string]string{
			"user_name": "lory",
			"user_id":   fmt.Sprintf("%d", 10000+i),
		}, 3)
		log.Debugf("send request [%v]", i)
		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(1 * time.Second)
	}
}

func NewGoMicroClient(typ gomicro.RegistryType) (c *gomicro.GoRPCClient) {
	var g *gomicro.GoRPC
	var endPoints []string
	g = gomicro.NewGoRPC(typ)
	switch typ {
	case gomicro.RegistryType_MDNS:
	case gomicro.RegistryType_ETCD:
		endPoints = strings.Split(END_POINTS_HTTP_ETCD, ",")
	case gomicro.RegistryType_CONSUL:
		endPoints = strings.Split(END_POINTS_HTTP_CONSUL, ",")
	case gomicro.RegistryType_ZOOKEEPER:
		endPoints = strings.Split(END_POINTS_TCP_ZOOKEEPER, ",")
	}
	return g.NewClient(endPoints...)
}
