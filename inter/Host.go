package inter

import ()

// Host for host maintain
type Host interface {
	Insert(host string, name string) error
	Query(name string) (string, error)
	Delete(name string) error
}