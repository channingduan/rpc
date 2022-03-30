package server

import (
	"context"
	"fmt"
	"github.com/channingduan/rpc/cache"
	"github.com/channingduan/rpc/config"
	"github.com/channingduan/rpc/database"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"time"
)

type RpcServer struct {
	server   *server.Server
	config   *config.Config
	cache    *cache.Cache
	database *database.Database
	methods  []config.Method
}

var RouterKey = "services:route"

func NewServer(config *config.Config) *RpcServer {

	return &RpcServer{
		server:   server.NewServer(),
		config:   config,
		cache:    cache.Register(&config.CacheConfig),
		database: database.Register(config),
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

	var members []string
	for _, method := range s.methods {
		members = append(members, fmt.Sprintf("%s.%s", s.config.ServicePath, method.Name))
		if err := s.server.RegisterFunctionName(s.config.ServicePath, method.Name, method.Func, ""); err != nil {
			log.Fatalf("register function cache error: %s", err)
		}
	}
	todo := context.TODO()
	c := s.cache.NewCache()
	val, _ := c.SScan(todo, RouterKey, 0, fmt.Sprintf("%s.*", s.config.ServicePath), 0).Val()
	c.SRem(todo, RouterKey, val)
	if err := c.SAdd(todo, RouterKey, members).Err(); err != nil {
		log.Fatalf("register function name cache error: %s", err)
	}
}

func (s *RpcServer) registryPlugin() {

	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", s.config.ServiceAddr),
		ConsulServers:  []string{s.config.RegistryConfig.Addr},
		BasePath:       s.config.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err := r.Start(); err != nil {
		log.Fatal("Register error: ", err)
	}

	s.server.Plugins.Add(r)
}
