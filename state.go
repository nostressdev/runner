package runner

type State interface {
	Ready() bool
	Alive() bool
}
