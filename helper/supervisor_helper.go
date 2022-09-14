package helper

import "github.com/mmaaskant/gophervisor/supervisor"

// StartSupervisor starts a supervisor.Supervisor instance based on the provided parameters.
func StartSupervisor(
	numberOfWorkers int,
	f func(p *supervisor.Publisher, d any, rch chan any),
) (*supervisor.Supervisor, *supervisor.Publisher, chan any) {
	sv := supervisor.NewSupervisor(numberOfWorkers)
	p, rch := sv.Register(f)
	return sv, p, rch
}
