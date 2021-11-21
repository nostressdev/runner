package stan

import "runner"

const (
	VariableGroup = "STAN"
	ClusterID     = "CLUSTER_ID"
	ClientID      = "CLIENT_ID"
	Url           = "URL"
)

type Config struct {
	StanConnURL   string
	StanClusterID string
	StanClientID  string
}

func NewConfigFromProvider(provider runner.VariableProvider) (*Config, error) {
	if err := provider.EnsureEnvironmentVariables(VariableGroup, Url, ClusterID, ClientID); err != nil {
		return nil, err
	}
	return &Config{
		StanConnURL:   provider.GetString(VariableGroup, Url),
		StanClusterID: provider.GetString(VariableGroup, ClusterID),
		StanClientID:  provider.GetString(VariableGroup, ClientID),
	}, nil
}
