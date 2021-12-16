package main

import (
	"fmt"
	"github.com/civet148/gomicro"
	"github.com/civet148/gomicro/example/echopb"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	"strings"
	"time"
)

const (
	SERVICE_NAME           = "echo"
	END_POINTS_HTTP_ETCD   = "http://127.0.0.1:2379"
	END_POINTS_HTTP_CONSUL = "http://127.0.0.1:8500"
	END_POINTS_ZOOKEEPER   = "127.0.0.1:2181"
)

func main() {

	c := NewGoMicroClient(gomicro.EndpointType_MDNS)
	service := echopb.NewEchoServerService(SERVICE_NAME, c)

	for i := 0; i < 10; i++ {
		ctx := gomicro.NewContext(map[string]string{
			"User_Name": "lory",
			"User_Id":   fmt.Sprintf("%d", 10000+i),
		}, 5)
		log.Debugf("send request [%v]", i)
		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(2 * time.Second)
	}
}

func NewGoMicroClient(typ gomicro.EndpointType) (c client.Client) {
	var g *gomicro.GoRPC
	var endPoints []string
	g = gomicro.NewGoRPC(typ)
	switch typ {
	case gomicro.EndpointType_MDNS:
	case gomicro.EndpointType_ETCD:
		endPoints = strings.Split(END_POINTS_HTTP_ETCD, ",")
	case gomicro.EndpointType_CONSUL:
		endPoints = strings.Split(END_POINTS_HTTP_CONSUL, ",")
	case gomicro.EndpointType_ZOOKEEPER:
		endPoints = strings.Split(END_POINTS_ZOOKEEPER, ",")
	}
	return g.NewClient(endPoints...)
}
