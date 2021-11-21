package stan

import (
	"context"
	"github.com/nats-io/stan.go"
	"runner"
	"time"
)

type Resource struct {
	stanConnURL   string
	stanClusterID string
	stanClientID  string
	StanConn      stan.Conn
}

func (res *Resource) Init(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return runner.NoDeadline
	}
	var err error
	res.StanConn, err = stan.Connect(res.stanClusterID, res.stanClientID, stan.NatsURL(res.stanConnURL), stan.ConnectWait(deadline.Sub(time.Now())))
	return err
}

func (res *Resource) Release(context.Context) error {
	res.StanConn.NatsConn().Close()
	return res.StanConn.Close()
}

func New(config *Config) *Resource {
	return &Resource{
		stanConnURL:   config.StanConnURL,
		stanClusterID: config.StanClusterID,
		stanClientID:  config.StanClientID,
	}
}
