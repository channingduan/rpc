package server

import (
	"fmt"
	"github.com/channingduan/rpc/config"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"time"
)

type RpcServer struct {
	server  *server.Server
	config  *config.Config
	methods []config.Method
}

func NewServer(config *config.Config) *RpcServer {

	return &RpcServer{
		server: server.NewServer(),
		config: config,
	}
}

func (s *RpcServer) Start() {
	s.registryPlugin()
	s.RegisterFunctionName()
	log.Infof("rpc server start: %v", s.config.ServiceAddr)
	if err := s.server.Serve("tcp", s.config.ServiceAddr); err != nil {
		panic(fmt.Sprintf("rpc server start error: %v", err))
	}
}

func (s *RpcServer) AddMethod(method config.Method) {
	s.methods = append(s.methods, method)
}

func (s *RpcServer) RegisterFunctionName() {
	for _, method := range s.methods {
		err := s.server.RegisterFunctionName(s.config.ServicePath, method.Name, method.Func, "")
		if err != nil {
			log.Fatal("RegisterFunctionName", err)
		}
	}

}

func (s *RpcServer) registryPlugin() {

	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", s.config.ServiceAddr),
		ConsulServers:  []string{s.config.RegistryAddr},
		BasePath:       s.config.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err := r.Start(); err != nil {
		log.Fatal("Register error: ", err)
	}

	s.server.Plugins.Add(r)
}
