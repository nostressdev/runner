package listener

import (
	"context"
	"fmt"
	"net"
)

type Resource struct {
	Listener   net.Listener
	listenAddr string
	needClose  bool
}

func (res *Resource) Init(ctx context.Context) error {
	var err error
	res.Listener, err = net.Listen("tcp", res.listenAddr)
	return err
}

func (res *Resource) Release(context.Context) error {
	if res.needClose {
		return res.Listener.Close()
	}
	return nil
}

func New(config *Config) *Resource {
	return &Resource{
		listenAddr: fmt.Sprintf("%s:%s", config.Addr, config.Port),
		needClose:  config.NeedClose,
	}
}
