package main

import (
	"fmt"
	"github.com/civet148/gomicro/v2"
	"github.com/civet148/gomicro/v2/example/echopb"
	"github.com/civet148/log"
	"time"
)

const (
	SERVICE_NAME         = "echo"
	END_POINTS_MDNS      = "mdns"
	END_POINTS_HTTP_ETCD = "etcd://192.168.20.108:2379"
)

func main() {

	c := gomicro.NewClient(END_POINTS_MDNS)
	service := echopb.NewEchoServerService(SERVICE_NAME, c)

	for i := 0; i < 20; i++ {
		ctx := gomicro.NewContext(map[string]string{
			"user_name": "lory",
			"user_id":   fmt.Sprintf("%d", 10000+i),
		})
		log.Debugf("send request [%v]", i)
		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(1 * time.Second)
	}
}
