package server

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/serverplugin"
	"time"
)

// 注册中心
func (s *RpcServer) pluginRegistry() {

	plugin := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s", s.config.ServiceAddr),
		ConsulServers:  []string{s.config.RegistryConfig.Addr},
		BasePath:       s.config.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	if err := plugin.Start(); err != nil {
		log.Fatal("Register error: ", err)
	}

	s.server.Plugins.Add(plugin)
}

// 限流
func (s *RpcServer) pluginRateLimit() {

	plugin := serverplugin.NewReqRateLimitingPlugin(time.Second, s.config.RateLimit, true)
	s.server.Plugins.Add(plugin)
}
