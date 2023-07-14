package main

import (
	"fmt"
	"github.com/civet148/gomicro/v2"
	"github.com/civet148/gomicro/v2/example/echopb"
	"github.com/civet148/gomicro/v2/strategy"
	"github.com/civet148/log"
	"time"
)

const (
	SERVICE_NAME         = "echo"
	END_POINTS_MDNS      = "mdns"
	END_POINTS_HTTP_ETCD = "etcd://192.168.2.9:2379"
)

func main() {
	log.SetLevel("debug")
	c := gomicro.NewClient(END_POINTS_HTTP_ETCD)
	service := echopb.NewEchoServerService(SERVICE_NAME, c)

	for i := 0; i < 9; i++ {
		ctx := gomicro.NewContext(map[string]string{
			"user_name": "lory",
			"user_id":   fmt.Sprintf("%d", 10000+i),
		})
		log.Debugf("send request [%v]", i)
		//opt := strategy.RoundRobin() //启用轮询策略(默认)
		//opt := strategy.Random() //启用随机策略
		//opt := strategy.Hash([]byte("192.168.2.100")) //启用哈希策略
		//opt := strategy.LeastLoad() //启用最小负载策略
		opt := strategy.Weight() //启用权重策略
		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}, opt); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(2 * time.Second)
	}
}
