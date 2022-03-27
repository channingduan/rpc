package server

import "github.com/channingduan/rpc/config"

type Server interface {
	AddMethod(method config.Method)
	RegisterFunctionName()
	Start()
}
