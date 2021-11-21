package runner

import "sync"

type State interface {
	Ready() bool
	Alive() bool
}

type stateImpl struct {
	sync.RWMutex
	alive bool
	ready bool
}

func (state *stateImpl) Ready() bool {
	state.RLock()
	defer state.RUnlock()
	return state.ready
}

func (state *stateImpl) Alive() bool {
	state.RLock()
	defer state.RUnlock()
	return state.alive
}

func (state *stateImpl) setAlive(alive bool) {
	state.Lock()
	defer state.Unlock()
	state.alive = alive
}

func (state *stateImpl) setReady(ready bool) {
	state.Lock()
	defer state.Unlock()
	state.ready = ready
}

func newImpl() State {
	return &stateImpl{
		alive: true,
		ready: false,
	}
}
