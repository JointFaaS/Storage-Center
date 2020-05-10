package inter

import "github.com/JointFaaS/Storage-Center/state"

// Host for host maintain
type Host interface {
	Insert(host string, name string) error
	Query(name string) (string, error)
	Delete(name string) error
	GetChan(name string) (chan state.InvalidEntry, error)
}
