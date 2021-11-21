package runner

import "context"

type Job interface {
	Run() error
	Shutdown(ctx context.Context) error
}
