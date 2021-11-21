package stan

import (
	"context"
	"github.com/nats-io/stan.go"
	"runner"
	"time"
)

const (
	VariableGroup = "STAN"
	ClusterID     = "CLUSTER_ID"
	ClientID      = "CLIENT_ID"
	Url           = "URL"
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

func New(provider runner.VariableProvider) (*Resource, error) {
	if err := provider.EnsureEnvironmentVariables(VariableGroup, Url, ClusterID, ClientID); err != nil {
		return nil, err
	}
	return &Resource{
		stanConnURL:   provider.GetString(VariableGroup, Url),
		stanClusterID: provider.GetString(VariableGroup, ClusterID),
		stanClientID:  provider.GetString(VariableGroup, ClientID),
	}, nil
}
