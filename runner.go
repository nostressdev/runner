package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	resources []Resource
	jobs      []Job
	config    Config
}

func New(config Config, resources []Resource, jobs []Job) *Runner {
	return &Runner{
		resources: resources,
		jobs:      jobs,
		config:    config,
	}
}

func (runner *Runner) init(ctx context.Context) error {
	initCtx, cancel := context.WithTimeout(ctx, runner.config.InitializationTimeout)
	defer cancel()
	errCh := make(chan error)
	go func(initCtx context.Context, errCh chan error) {
		for _, resource := range runner.resources {
			if err := resource.Init(initCtx); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}(initCtx, errCh)
	select {
	case <-time.After(runner.config.InitializationTimeout):
		return InitializationDeadline
	case err := <-errCh:
		return err
	}
}

func (runner *Runner) Run() error {
	ctx := context.Background()

	if err := runner.init(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	jobsDoneCh := make(chan error)

	for _, job := range runner.jobs {
		wg.Add(1)
		go func(jobsDoneCh chan error, job Job) {
			if err := job.Run(); err != nil {
				jobsDoneCh <- err
			}
			wg.Done()
		}(jobsDoneCh, job)
	}

	go func(jobsDoneCh chan error) {
		wg.Wait()
		jobsDoneCh <- nil
	}(jobsDoneCh)

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case jobsError := <-jobsDoneCh:
		shutdownError := runner.Shutdown()
		return NewMultiError(jobsError, shutdownError)
	case <-sig:
		return runner.Shutdown()
	case <-ctx.Done():
		return runner.Shutdown()
	}
}

func (runner *Runner) Shutdown() error {

	terminationCtx, cancel := context.WithTimeout(context.Background(), runner.config.TerminationTimeout)
	defer cancel()

	errCh := make(chan error)

	go func(terminationCtx context.Context, errCh chan error) {
		terminationErrors := make([]error, 0)
		// Stop running jobs
		for _, job := range runner.jobs {
			if err := job.Shutdown(terminationCtx); err != nil {
				terminationErrors = append(terminationErrors, err)
			}
		}
		// Release resources
		for _, resource := range runner.resources {
			if err := resource.Release(terminationCtx); err != nil {
				terminationErrors = append(terminationErrors, err)
			}
		}
		errCh <- NewMultiError(terminationErrors...)
	}(terminationCtx, errCh)

	select {
	case <-time.After(runner.config.TerminationTimeout):
		return ShutdownDeadline
	case err := <-errCh:
		return err
	}

}
