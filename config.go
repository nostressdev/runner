package runner

import "time"

type Config struct {
	InitializationTimeout time.Duration
	TerminationTimeout    time.Duration
}

func DefaultConfig() Config {
	return Config{
		InitializationTimeout: time.Minute,
		TerminationTimeout:    time.Minute,
	}
}