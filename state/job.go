package state

import (
	"context"
	"errors"
	"github.com/nostressdev/runner"
	"github.com/nostressdev/runner/listener"
	"net/http"
)

type readyLiveHttpJob struct {
	listener *listener.Resource
	finished chan error
	state    runner.State
	server   *http.Server
}

func (job *readyLiveHttpJob) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		if job.state.Alive() {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
	})
	mux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		if job.state.Ready() {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
	})

	job.server = &http.Server{
		Handler: mux,
	}

	if err := job.server.Serve(job.listener.Listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		job.finished <- err
		return err
	}
	return nil
}

func (job *readyLiveHttpJob) Shutdown(ctx context.Context) error {
	go func() {
		if err := job.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			job.finished <- err
		} else {
			job.finished <- nil
		}
	}()
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-job.finished:
		return err
	}
}

func NewReadyLiveHttpJob(listener *listener.Resource, state runner.State) runner.Job {
	return &readyLiveHttpJob{
		listener: listener,
		finished: make(chan error, 1),
		state:    state,
	}
}
