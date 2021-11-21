package runner

import (
	"context"
)

type Resource interface {
	Init(ctx context.Context) error
	Release(ctx context.Context) error
}
