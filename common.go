package gomicro

import (
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/registry/mdns"
	sgrpc "github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
	"strings"
	"time"
)

/*
	"etcd://192.168.1.108:2379,192.168.1.109:2379"
	"consul://192.168.1.108:8500,192.168.1.109:8500"
	"zk://192.168.1.108:2181, 192.168.1.109:2181"
*/
func ParseRegistry(strRegistry string) (typ RegistryType, endpoints []string) {
	var strAddress string
	//--registry "etcd://127.0.0.1:2379,127.0.0.1:2380"
	ss := strings.Split(strRegistry, "://")
	count := len(ss)
	if count > 1 {
		strRegName := strings.ToLower(ss[0])
		strAddress = strings.ToLower(ss[1])
		switch strRegName {
		case "etcd":
			typ = RegistryType_ETCD
		//case "consul":
		//	typ = RegistryType_CONSUL
		//case "zookeeper", "zk":
		//	typ = RegistryType_ZOOKEEPER
		default:
			log.Warnf("Unknown registry name [%s], use default MDNS", strRegName)
			typ = RegistryType_MDNS
		}
	} else {
		typ = RegistryType_MDNS
	}
	log.Infof("registry type [%s] address %+v", typ.String(), strAddress)
	endpoints = strings.Split(strAddress, ",")
	return typ, endpoints
}

func newRegistry(registryType RegistryType, endPoints ...string) (r registry.Registry) {
	var opts []registry.Option
	opts = append(opts, registry.Addrs(endPoints...))

	switch registryType {
	case RegistryType_MDNS:
		r = mdns.NewRegistry()
	case RegistryType_ETCD:
		r = etcd.NewRegistry(opts...)
	//case RegistryType_CONSUL:
	//	r = consul.NewRegistry(opts...)
	//case RegistryType_ZOOKEEPER:
	//	r = zookeeper.NewRegistry(opts...)
	default:
		log.Panic("end point type [%+v] illegal", registryType)
	}
	log.Debugf("[%+v] end points %+v -> registry [%+v]", registryType, endPoints, r)
	return
}

func newOptions(registryType RegistryType,  discovery *Discovery, reg registry.Registry) []service.Option {
	var options []service.Option
	if reg == nil {
		log.Panic("[%+v] discovery [%+v] -> registry is nil", registryType, discovery)
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
		log.Panic("discover service name is nil")
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
