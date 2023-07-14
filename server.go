package gomicro

import (
	"github.com/civet148/log"
	"github.com/micro/go-micro/v2/server"
)

func (s *GoRPCServer) Close() error {
	services, err := s.registry.GetService(s.discovery.ServiceName)
	if err != nil {
		return log.Errorf("registry get service by service name %s error [%s]", s.discovery.ServiceName, err)
	}
	for _, v := range services {
		err = s.registry.Deregister(v)
		if err != nil {
			return log.Errorf("deregister service by service name %s error [%s]", s.discovery.ServiceName, err)
		}
	}
	if err := s.server.Stop(); err != nil {
		return log.Errorf("server stop error [%s]", err)
	}
	return nil
}

// Initialise options
func (s *GoRPCServer) Init(opts ...server.Option) error {
	return s.server.Init(opts...)
}

// Retrieve the options
func (s *GoRPCServer) Options() server.Options {
	return s.server.Options()
}

// Register a handler
func (s *GoRPCServer) Handle(h server.Handler) error {
	return s.server.Handle(h)
}

// Create a new handler
func (s *GoRPCServer) NewHandler(h interface{}, opts ...server.HandlerOption) server.Handler {
	return s.server.NewHandler(h, opts...)
}

// Create a new subscriber
func (s *GoRPCServer) NewSubscriber(topic string, sb interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return s.server.NewSubscriber(topic, sb, opts...)
}

// Register a subscriber
func (s *GoRPCServer) Subscribe(subscriber server.Subscriber) error {
	return s.server.Subscribe(subscriber)
}

// Start the server
func (s *GoRPCServer) Start() error {
	return s.server.Start()
}

// Stop the server
func (s *GoRPCServer) Stop() error {

	return s.server.Stop()
}

// Server implementation
func (s *GoRPCServer) String() string {

	return s.server.String()
}

// Server weight of load balance
func (s *GoRPCServer) Weight() int {
	return s.weight
}
