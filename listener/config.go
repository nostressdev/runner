package listener

import "github.com/nostressdev/runner"

const (
	Port = "PORT"
	Addr = "ADDR"
)

type Config struct {
	Addr string
	Port string
}

func NewConfigFromProvider(provider runner.VariableProvider, group string) (*Config, error) {
	if err := provider.EnsureEnvironmentVariables(group, Port, Addr); err != nil {
		return nil, err
	}
	return &Config{
		Addr: provider.GetString(group, Addr),
		Port: provider.GetString(group, Port),
	}, nil
}
