package gomicro

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/registry/mdns"
	sgrpc "github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/registry/zookeeper/v2"
	"time"
)

func newRegistry(registryType RegistryType, endPoints ...string) (r registry.Registry) {
	var opts []registry.Option
	opts = append(opts, registry.Addrs(endPoints...))

	switch registryType {
	case RegistryType_MDNS:
		r = mdns.NewRegistry()
	case RegistryType_ETCD:
		r = etcd.NewRegistry(opts...)
	case RegistryType_CONSUL:
		r = consul.NewRegistry(opts...)
	case RegistryType_ZOOKEEPER:
		r = zookeeper.NewRegistry(opts...)
	default:
		panic(fmt.Errorf("end point type [%+v] illegal", registryType))
	}
	log.Debugf("[%+v] end points %+v -> registry [%+v]", registryType, endPoints, r)
	return
}

func newOptions(registryType RegistryType,  discovery *Discovery, reg registry.Registry) []service.Option {
	var options []service.Option
	if reg == nil {
		panic(fmt.Errorf("[%+v] discovery [%+v] -> registry is nil", registryType, discovery))
	}
	options = append(options, service.Registry(reg))
	options = append(options, service.RegisterInterval(time.Duration(discovery.Interval)*time.Second))
	options = append(options, service.RegisterTTL(time.Duration(discovery.TTL)*time.Second))
	options = append(options, service.Name(discovery.ServiceName))
	options = append(options, service.Address(discovery.RpcAddr))
	return options
}


//NewServer new a go-micro server
func newRpcServer(registryType RegistryType, discovery *Discovery, maxMsgSize int) (s *GoRPCServer) { // returns go-micro server object
	log.Debugf("endpoint type [%v] discovery [%+v]", registryType, discovery)
	if len(discovery.Endpoints) == 0 {
		registryType = RegistryType_MDNS
	}
	if discovery.ServiceName == "" {
		panic("discover service name is nil")
	}
	if discovery.Interval == 0 {
		discovery.Interval = DISCOVERY_DEFAULT_INTERVAL
	}
	if discovery.TTL == 0 {
		discovery.TTL = DISCOVERY_DEFAULT_TTL
	}
	s = &GoRPCServer{}
	var options []service.Option
	reg := newRegistry(registryType, discovery.Endpoints...)
	options = newOptions(registryType, discovery, reg)
	rpc := grpc.NewService(options...)
	opt := sgrpc.MaxMsgSize(maxMsgSize)

	s.server = rpc.Server()
	s.registry = reg
	s.discovery = discovery
	s.registryType = registryType
	s.options = options
	s.maxMsgSize = maxMsgSize
	if err := s.server.Init(opt); err != nil {
		log.Panic("initialize server option error [%s]", err)
	}
	return s
}