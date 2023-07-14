package gomicro

import (
	"context"
	"fmt"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	cgrpc "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
	"time"
)

var (
	DefaultMaxMsgSize = 512 * 1024 * 1024 //default allow 512 MiB data transport
)

const (
	DISCOVERY_DEFAULT_INTERVAL = 3
	DISCOVERY_DEFAULT_TTL      = 10
	DEFAULT_RPC_TIMEOUT        = 30
	LOAD_BALANCE_KEY_WEIGTHT   = "gomicro.LoadBalanceWeight"
)

type RegistryType int

const (
	RegistryType_MDNS RegistryType = 0 // multicast DNS
	RegistryType_ETCD RegistryType = 1 // etcd
	//RegistryType_CONSUL    RegistryType = 2 // consul
	//RegistryType_ZOOKEEPER RegistryType = 3 // zookeeper
)

func (t RegistryType) String() string {
	switch t {
	case RegistryType_MDNS:
		return "RegistryType_MDNS"
	case RegistryType_ETCD:
		return "RegistryType_ETCD"
		//case RegistryType_CONSUL:
		//	return "RegistryType_CONSUL"
		//case RegistryType_ZOOKEEPER:
		//	return "RegistryType_ZOOKEEPER"
	}
	return "RegistryType_Unknown"
}

type ServerOption struct {
	ServiceName string            // register service name [required]
	RpcAddr     string            // register server RPC address [required]
	Interval    int               // register interval default 3 seconds [optional]
	TTL         int               // register TTL default 10 seconds [optional]
	Weight      int               // server weight of load balance [optional]
	Metadata    map[string]string // register node metadata [optional]
}

type discovery struct {
	ServiceName string   // register service name [required]
	RpcAddr     string   // register server RPC address [required]
	Interval    int      // register interval default 3 seconds [optional]
	TTL         int      // register TTL default 10 seconds [optional]
	Weight      int      // server weight of load balance [optional]
	Endpoints   []string // register endpoints of etcd/consul/zookeeper eg. ["192.168.0.10:2379","192.168.0.11:2379"]
}

type GoRPC struct {
	maxMsgSize   int
	registryType RegistryType //end point type
}

type GoRPCClient struct {
	registry registry.Registry
	client   client.Client
}

type GoRPCServer struct {
	weight       int
	registry     registry.Registry
	server       server.Server
	maxMsgSize   int
	discovery    *discovery
	registryType RegistryType
	options      []service.Option
}

func init() {

}

func NewClient(strRegistry string) (c *GoRPCClient) {
	var g *GoRPC
	var endPoints []string
	var typ RegistryType
	typ, endPoints = parseRegistry(strRegistry)
	g = newGoRPC(typ)
	return g.NewClient(endPoints...)
}

func NewServer(strRegistry string, option *ServerOption) (s *GoRPCServer) {
	var g *GoRPC
	var endPoints []string
	var typ RegistryType
	typ, endPoints = parseRegistry(strRegistry)
	g = newGoRPC(typ)
	return g.NewServer(&discovery{
		ServiceName: option.ServiceName,
		RpcAddr:     option.RpcAddr,
		Interval:    option.Interval,
		Weight:      option.Weight,
		TTL:         option.TTL,
		Endpoints:   endPoints,
	}, option.Metadata)
}

// sizes: max send or receive msg size in byte
func newGoRPC(registryType RegistryType, maxSize ...int) (g *GoRPC) {
	maxMsgSize := DefaultMaxMsgSize
	if len(maxSize) != 0 {
		maxMsgSize = maxSize[0]
	}
	return &GoRPC{
		registryType: registryType,
		maxMsgSize:   maxMsgSize,
	}
}

// NewContext
// md -> metadata of RPC call, set to nil if have no any meta-data
func NewContext(md map[string]string) context.Context {
	return metadata.NewContext(context.Background(), md)
}

// NewContextWithTimeout
// md -> metadata of RPC call, set to nil if have no any meta-data
// timeout -> timeout seconds of RPC call
func NewContextWithTimeout(md map[string]string, timeout int) (context.Context, context.CancelFunc) {
	ctx, closer := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return metadata.NewContext(ctx, md), closer
}

// FromContext get metadata from context
func FromContext(ctx context.Context) (md metadata.Metadata) {
	md, _ = metadata.FromContext(ctx)
	return
}

// NewClient new a go-micro client
func (g *GoRPC) NewClient(endPoints ...string) (c *GoRPCClient) { // returns go-micro client object
	var options []service.Option
	log.Debugf("endpoint type [%v] end points [%+v]", g.registryType, endPoints)
	optSend := cgrpc.MaxSendMsgSize(g.maxMsgSize)
	optRecv := cgrpc.MaxRecvMsgSize(g.maxMsgSize)

	reg := newRegistry(g.registryType, endPoints...)
	if reg != nil {
		options = append(options, service.Registry(reg))
	}

	rpc := grpc.NewService(options...)
	c = &GoRPCClient{
		registry: reg,
		client:   rpc.Client(),
	}
	if err := c.client.Init(optSend, optRecv); err != nil {
		log.Panic("initialize client option error [%s]", err)
	}
	return c
}

// NewServer new a go-micro server with metadata
func (g *GoRPC) NewServer(discovery *discovery, metadata map[string]string) (s *GoRPCServer) { // returns go-micro server object
	if discovery.Weight > 0 {
		if metadata == nil {
			metadata = make(map[string]string)
		}
		metadata[LOAD_BALANCE_KEY_WEIGTHT] = fmt.Sprintf("%v", discovery.Weight)
	}
	return newRpcServer(g.registryType, discovery, g.maxMsgSize, metadata)
}
