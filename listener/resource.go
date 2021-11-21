package listener

import (
	"context"
	"net"
)

type Resource struct {
	listenAddr string
	Listener   net.Listener
}

func (res *Resource) Init(ctx context.Context) error {
	var err error
	res.Listener, err = net.Listen("tcp", res.listenAddr)
	return err
}

func (res *Resource) Release(context.Context) error {
	return res.Listener.Close()
}

func New(listenAddr string) *Resource {
	return &Resource{
		listenAddr: listenAddr,
	}
}
