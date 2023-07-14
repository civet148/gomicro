package strategy

//algorithms refer to https://github.com/zehuamama/balancer
import (
	"github.com/civet148/gomicro/v2"
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"hash/crc32"
	"math/rand"
	"strconv"
	"sync"
)

type weightNode struct {
	Weight    int            //固定权重
	CurWeight int            //当前权重
	Node      *registry.Node //服务节点
}

var idx uint64
var effectiveWeight int
var locker sync.RWMutex
var loads = make(map[string]uint64)
var weights = make(map[string]*weightNode)

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

// Weight 权重策略(weight strategy)
func Weight() client.CallOption {

	return client.WithSelectOption(selector.WithStrategy(func(services []*registry.Service) selector.Next {
		var nodes []*registry.Node
		for _, s := range services {
			nodes = append(nodes, s.Nodes...)
		}
		return func() (node *registry.Node, err error) {
			if len(nodes) == 0 {
				log.Errorf("no available nodes")
				resetWeightNodes()
				return nil, selector.ErrNoneAvailable
			}
			return getExpectedNode(nodes), nil
		}
	}))
}

func resetWeightNodes() {
	locker.Lock()
	defer locker.Unlock()
	effectiveWeight = 0
	for k := range weights {
		delete(weights, k)
	}
}

func getExpectedNode(nodes []*registry.Node) (node *registry.Node) {
	var reset bool
	var dic = make(map[string]*registry.Node)
	for _, n := range nodes {
		dic[n.Address] = n
	}

	locker.Lock()
	defer locker.Unlock()

	//清理已下线服务节点信息
	for _, w := range weights {
		if _, ok := dic[w.Node.Address]; !ok {
			reset = true
			delete(weights, w.Node.Address)
		}
	}
	//更新本地节点信息和权重
	for _, n := range nodes {
		strVal := n.Metadata[gomicro.LOAD_BALANCE_KEY_WEIGTHT]
		weight, _ := strconv.Atoi(strVal)
		if w, ok := weights[n.Address]; !ok {
			reset = true
			weights[n.Address] = &weightNode{
				Weight:    weight,
				CurWeight: 0,
				Node:      n,
			}
		} else {
			w.Node = n
			if w.Weight != weight {
				reset = true
				w.CurWeight = 0
				w.Weight = weight
			}
		}
	}
	if reset {
		//重新计算有效权重值
		effectiveWeight = 0
		for _, w := range weights {
			w.CurWeight = 0
			effectiveWeight += w.Weight
		}
	}
	var weight int
	var wgt *weightNode
	for _, w := range weights {
		w.CurWeight += w.Weight
		if weight == 0 || w.CurWeight > weight {
			wgt = w
			weight = w.CurWeight
		}
	}
	wgt.CurWeight -= effectiveWeight
	return wgt.Node
}
