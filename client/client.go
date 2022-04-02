package client

import (
	"context"
	"fmt"
	"github.com/channingduan/rpc/config"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type RpcClient struct {
	config     *config.Config
	clients    map[string]client.XClient
	discovery  client.ServiceDiscovery
	failMode   client.FailMode
	selectMode client.SelectMode
	option     client.Option
}

// NewClient 初始化客户端
func NewClient(config *config.Config) *RpcClient {

	discovery, err := client.NewConsulDiscovery(config.BasePath, config.ServicePath, []string{config.RegistryConfig.Addr}, nil)
	if err != nil {
		panic(fmt.Sprintf("server discovery error: %v", err))
	}
	option := client.DefaultOption
	option.SerializeType = protocol.JSON
	return &RpcClient{
		config:     config,
		clients:    make(map[string]client.XClient),
		discovery:  discovery,
		option:     option,
		failMode:   client.Failtry,
		selectMode: client.RandomSelect,
	}
}

func (c *RpcClient) getClient(serverPath string) (client.XClient, error) {

	if c.clients[serverPath] == nil {
		c.clients[serverPath] = client.NewXClient(c.config.ServicePath, c.failMode, c.selectMode, c.discovery, c.option)
	}

	return c.clients[serverPath], nil
}

// Call RPC 方法调用
func (c *RpcClient) Call(ctx context.Context, serverPath, method string, request config.Request) (*config.Response, error) {

	fc, err := c.getClient(serverPath)
	if err != nil {
		return nil, err
	}
	var response config.Response
	if err := fc.Call(ctx, method, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
