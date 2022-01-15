package gomicro

import (
	"context"
	"fmt"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	cgrpc "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/registry/mdns"
	"github.com/micro/go-micro/v2/server"
	sgrpc "github.com/micro/go-micro/v2/server/grpc"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/registry/zookeeper/v2"
	"time"
)

var (
	DefaultMaxMsgSize = 512 * 1024 * 1024 //default allow 512 MiB data transport
)

const (
	DISCOVERY_DEFAULT_INTERVAL = 3
	DISCOVERY_DEFAULT_TTL      = 10
	DEFAULT_RPC_TIMEOUT        = 30
)

type RegistryType int

const (
	RegistryType_MDNS      RegistryType = 0 // multicast DNS
	RegistryType_ETCD      RegistryType = 1 // etcd
	RegistryType_CONSUL    RegistryType = 2 // consul
	RegistryType_ZOOKEEPER RegistryType = 3 // zookeeper
)

func (t RegistryType) String() string {
	switch t {
	case RegistryType_MDNS:
		return "RegistryType_MDNS"
	case RegistryType_ETCD:
		return "RegistryType_ETCD"
	case RegistryType_CONSUL:
		return "RegistryType_CONSUL"
	case RegistryType_ZOOKEEPER:
		return "RegistryType_ZOOKEEPER"
	}
	return "RegistryType_Unknown"
}

type Discovery struct {
	ServiceName string   // register service name [required]
	RpcAddr     string   // register server RPC address [required]
	Interval    int      // register interval default 3 seconds [optional]
	TTL         int      // register TTL default 10 seconds [optional]
	Endpoints   []string // register endpoints of etcd/consul/zookeeper eg. ["192.168.0.10:2379","192.168.0.11:2379"]
}

type GoRPC struct {
	maxMsgSize   int
	registryType RegistryType //end point type
}

func init() {

}

//sizes: max send or receive msg size in byte
func NewGoRPC(registryType RegistryType, maxSize ...int) (g *GoRPC) {
	maxMsgSize := DefaultMaxMsgSize
	if len(maxSize) != 0 {
		maxMsgSize = maxSize[0]
	}
	return &GoRPC{
		registryType: registryType,
		maxMsgSize:   maxMsgSize,
	}
}

//NewContext
//md -> metadata of RPC call, set to nil if have no any meta-data
//timeout -> timeout seconds of RPC call, if <=0 will set it to DEFAULT_RPC_TIMEOUT
func NewContext(md map[string]string, timeout int) context.Context {
	var ctx = context.Background()
	if timeout <= 0 {
		timeout = DEFAULT_RPC_TIMEOUT
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	return metadata.NewContext(ctx, md)
}

//FromContext get metadata from context
func FromContext(ctx context.Context) (md metadata.Metadata) {
	md, _ = metadata.FromContext(ctx)
	return
}

//NewClient new a go-micro client
func (g *GoRPC) NewClient(endPoints ...string) (c client.Client) { // returns go-micro client object
	var options []service.Option

	log.Debugf("endpoint type [%v] end points [%+v]", g.registryType, endPoints)
	optSend := cgrpc.MaxSendMsgSize(g.maxMsgSize)
	optRecv := cgrpc.MaxRecvMsgSize(g.maxMsgSize)

	reg := g.newRegistry(endPoints...)
	if reg != nil {
		options = append(options, service.Registry(reg))
	}

	rpc := grpc.NewService(options...)
	c = rpc.Client()
	if err := c.Init(optSend, optRecv); err != nil {
		log.Panic("initialize client option error [%s]", err)
	}
	return c
}

//NewServer new a go-micro server
func (g *GoRPC) NewServer(discovery *Discovery) (s server.Server) { // returns go-micro server object
	log.Debugf("endpoint type [%v] discovery [%+v]", g.registryType, discovery)
	if len(discovery.Endpoints) == 0 {
		g.registryType = RegistryType_MDNS
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

	reg := g.newRegistry(discovery.Endpoints...)

	var options []service.Option
	if reg == nil {
		panic(fmt.Errorf("[%+v] discovery [%+v] -> registry is nil", g.registryType, discovery))
	}

	options = append(options, service.Registry(reg))
	options = append(options, service.RegisterInterval(time.Duration(discovery.Interval)*time.Second))
	options = append(options, service.RegisterTTL(time.Duration(discovery.TTL)*time.Second))
	options = append(options, service.Name(discovery.ServiceName))
	options = append(options, service.Address(discovery.RpcAddr))
	rpc := grpc.NewService(options...)
	opt := sgrpc.MaxMsgSize(g.maxMsgSize)

	s = rpc.Server()
	if err := s.Init(opt); err != nil {
		log.Panic("initialize server option error [%s]", err)
	}
	return s
}

func (g *GoRPC) newRegistry(endPoints ...string) (r registry.Registry) {
	var opts []registry.Option
	opts = append(opts, registry.Addrs(endPoints...))

	switch g.registryType {
	case RegistryType_MDNS:
		r = mdns.NewRegistry()
	case RegistryType_ETCD:
		r = etcd.NewRegistry(opts...)
	case RegistryType_CONSUL:
		r = consul.NewRegistry(opts...)
	case RegistryType_ZOOKEEPER:
		r = zookeeper.NewRegistry(opts...)
	default:
		panic(fmt.Errorf("end point type [%+v] illegal", g.registryType))
	}
	log.Debugf("[%+v] end points [%+v] -> registry [%+v]", g.registryType, endPoints, r)
	return
}
