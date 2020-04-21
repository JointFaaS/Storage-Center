package inter

import ()

// State for maintain StateInterface
type State interface {
	ChangeStatus(token string, host string) (string , string, error)
	Query(token string) (string, error)
}