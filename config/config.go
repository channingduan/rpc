package config

import "context"

// Config 配置文件
type Config struct {
	BasePath     string
	ServicePath  string
	ServiceName  string
	ServiceAddr  string
	RegistryAddr string
}

type Method struct {
	Name string
	Func func(ctx context.Context, res *Request, req *Response) error
}
