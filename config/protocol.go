package config

import "context"

// Request 请求解析
type Request struct {
	Message string `json:"message"`
}

// Response 返回协议
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HttpResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Config struct {
	BasePath       string                    `json:"base_path"`
	ServicePath    string                    `json:"service_path"`
	ServiceName    string                    `json:"service_name"`
	ServiceAddr    string                    `json:"service_addr"`
	RegistryConfig RegistryConfig            `json:"registry_config"`
	RateLimit      int64                     `json:"rate_limit"`
	DatabaseConfig map[string]DatabaseConfig `json:"database_config"`
	CacheConfig    CacheConfig               `json:"cache_config"`
	TokenConfig    TokenConfig               `json:"token_config"`
}

// TokenConfig config
type TokenConfig struct {
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
	AccessExpire  int    `json:"access_expire"`
	RefreshExpire int    `json:"refresh_expire"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string            `json:"driver"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Database string            `json:"database"`
	Sources  []DatabaseConnect `json:"sources"`
	Replicas []DatabaseConnect `json:"replicas"`
}

type DatabaseConnect struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// RegistryConfig 服务发现配置
type RegistryConfig struct {
	Driver string `json:"driver"`
	Addr   string `json:"addr"`
}

type CacheConfig struct {
	Driver   string `json:"driver"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Method 服务结构
type Method struct {
	Name   string                                                       `json:"name"`
	Router string                                                       `json:"router"`
	Func   func(ctx context.Context, res *Request, req *Response) error `json:"func"`
}
