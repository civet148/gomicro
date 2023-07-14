package strategy

//algorithms refer to https://github.com/zehuamama/balancer
import (
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"hash/crc32"
	"math/rand"
	"sync"
)

var idx uint64
var locker sync.RWMutex
var loads = make(map[string]uint64)

// Random 随机策略(random)
func Random() client.CallOption {

	return client.WithSelectOption(selector.WithStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node
		for _, s := range services {
			nodes = append(nodes, s.Nodes...)
		}

		return func() (*registry.Node, error) {
			if len(nodes) == 0 {
				log.Errorf("no available nodes")
				return nil, selector.ErrNoneAvailable
			}
			var i = rand.Int()
			node := nodes[i%len(nodes)]
			return node, nil
		}
	}))
}

// RoundRobin 轮询策略(round-robin)
func RoundRobin() client.CallOption {

	return client.WithSelectOption(selector.WithStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node
		for _, s := range services {
			nodes = append(nodes, s.Nodes...)
		}

		return func() (*registry.Node, error) {
			if len(nodes) == 0 {
				log.Errorf("no available nodes")
				return nil, selector.ErrNoneAvailable
			}
			locker.Lock()
			node := nodes[idx%uint64(len(nodes))]
			idx++
			locker.Unlock()
			return node, nil
		}
	}))
}

// Hash 哈希询策略(hash)
func Hash(key []byte) client.CallOption {

	return client.WithSelectOption(selector.WithStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node
		for _, s := range services {
			nodes = append(nodes, s.Nodes...)
		}
		return func() (node *registry.Node, err error) {
			if len(nodes) == 0 {
				log.Errorf("no available nodes")
				return nil, selector.ErrNoneAvailable
			}
			i := crc32.ChecksumIEEE(key) % uint32(len(nodes))
			node = nodes[i]
			return node, nil
		}
	}))
}

// LeastLoad 最小负载策略(least-load)
func LeastLoad() client.CallOption {

	return client.WithSelectOption(selector.WithStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node
		for _, s := range services {
			nodes = append(nodes, s.Nodes...)
		}
		return func() (node *registry.Node, err error) {
			if len(nodes) == 0 {
				log.Errorf("no available nodes")
				return nil, selector.ErrNoneAvailable
			}
			var min uint64
			locker.Lock()
			for _, n := range nodes {
				count, ok := loads[n.Id]
				if !ok || count == 0 {
					node = n
					min = count
					break
				}
				if min == 0 || count < min {
					node = n
					min = count
				}
			}
			loads[node.Id] = min + 1
			locker.Unlock()
			return node, nil
		}
	}))
}
