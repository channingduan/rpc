package service

import (
	"context"
	"fmt"
	"github.com/channingduan/rpc/config"
	"github.com/oscto/ky3k"
	"github.com/smallnest/rpcx/log"
	"time"
)

type Greeter struct{}

func (g *Greeter) Say(ctx context.Context, name *string, reply *string) error {

	*reply = fmt.Sprintf("Hello %s!", *name)
	return nil
}

func (g *Greeter) Hello(ctx context.Context, resp interface{}, reply *config.Response) error {

	message := "hello world"
	reply.Message = message
	log.Debug("Hello...", time.Now().String())
	return nil
}

type Human struct {
}

func (h *Human) Hello(ctx context.Context, resp *config.Request, reply *config.Response) error {

	message := "Human Hello"
	reply.Message = message
	log.Debug("Human Hello...", time.Now().String())

	return nil
}

func Jerry(ctx context.Context, res *config.Request, req *config.Response) error {

	req.Message = "hello jerry"
	fmt.Println("jerry:")
	fmt.Println("res", ky3k.JsonToString(res))
	fmt.Println("req", ky3k.JsonToString(req))
	return nil
}
