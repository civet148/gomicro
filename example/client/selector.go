package main

import (
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

func SelectorOption() client.CallOption {
	return client.WithSelectOption(selector.WithFilter(func(services []*registry.Service) []*registry.Service {
		for i, s := range services {
			if len(s.Nodes) == 0 {
				log.Warnf("selector service[%d] name [%s] nodes is 0", i, s.Name)
				return services
			}
			log.Json("service nodes", s.Nodes)
		}
		return services
	}))
}
