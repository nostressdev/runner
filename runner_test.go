package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	ErrFailInit    = errors.New("fail init")
	ErrFailRelease = errors.New("fail release")
)

type TestResource struct {
	InitFailTime    time.Duration
	ReleaseFailTime time.Duration
}

func (f *TestResource) Init(ctx context.Context) error {
	if f.InitFailTime != 0 {
		<-time.After(f.InitFailTime)
		return ErrFailInit
	}
	return nil
}

func (f *TestResource) Release(ctx context.Context) error {
	if f.ReleaseFailTime != 0 {
		select {
		case <-time.After(f.ReleaseFailTime):
			return ErrFailRelease
		case <-ctx.Done():
			return context.DeadlineExceeded
		}
	}
	return nil
}

type TestJob struct {
	WorkTime      time.Duration
	WorkError     error
	ShutdownTime  time.Duration
	ShutdownError error
	ShutdownChan  chan struct{}
	FinishChan    chan struct{}
}

func (tj *TestJob) Run() error {
	defer func() {
		tj.FinishChan <- struct {
		}{}
	}()
	select {
	case <-time.After(tj.WorkTime):
		return tj.WorkError
	case <-tj.ShutdownChan:
		return nil
	}
}

func (tj *TestJob) Shutdown(ctx context.Context) error {
	defer func() {
		tj.ShutdownChan <- struct{}{}
	}()
	select {
	case <-time.After(tj.ShutdownTime):
		return tj.ShutdownError
	case <-ctx.Done():
		return context.DeadlineExceeded
	}
}

func Test_OK(t *testing.T) {
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, nil, nil)
	assert.NoError(t, runner.Run())
}

func Test_FailInit(t *testing.T) {
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, []Resource{
		&TestResource{
			InitFailTime: time.Millisecond,
		},
	}, nil)
	assert.EqualError(t, runner.Run(), ErrFailInit.Error())
}

func Test_FailRelease(t *testing.T) {
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, []Resource{
		&TestResource{
			ReleaseFailTime: time.Millisecond,
		},
	}, nil)
	assert.Error(t, runner.Run(), ErrFailRelease.Error())
}

func Test_RunJob(t *testing.T) {
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, []Resource{
		&TestResource{},
		&TestResource{},
	}, []Job{
		&TestJob{
			ShutdownTime: time.Millisecond * 500,
			FinishChan:   make(chan struct{}, 1),
			ShutdownChan: make(chan struct{}, 1),
		},
	})
	assert.NoError(t, runner.Run())
}

func Test_RunJobFailWork(t *testing.T) {
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, []Resource{
		&TestResource{},
		&TestResource{},
	}, []Job{
		&TestJob{
			WorkTime:     500,
			ShutdownTime: time.Millisecond * 500,
			FinishChan:   make(chan struct{}, 1),
			ShutdownChan: make(chan struct{}, 1),
		},
	})
	assert.NoError(t, runner.Run())
}

func Test_RunJobFailCheckOtherStop(t *testing.T) {
	err := fmt.Errorf("j1 fail")
	j1 := &TestJob{
		WorkTime:     time.Millisecond * 100,
		WorkError:    err,
		FinishChan:   make(chan struct{}, 1),
		ShutdownChan: make(chan struct{}, 1),
	}
	j2 := &TestJob{
		WorkTime:     time.Second * 5000,
		FinishChan:   make(chan struct{}, 1),
		ShutdownChan: make(chan struct{}, 1),
	}
	runner := New(Config{
		InitializationTimeout: time.Second,
		TerminationTimeout:    time.Second,
	}, []Resource{
		&TestResource{},
		&TestResource{},
	}, []Job{
		j1,
		j2,
	})
	go func() {
		assert.Error(t, runner.Run(), err.Error())
	}()
	for i := 0; i < 2; i += 1 {
		select {
		case <-time.After(time.Second * 1):
			t.Fatalf("expected stop job")
		case <-j2.FinishChan:
			continue
		case <-j1.FinishChan:
			continue
		}
	}
}
