package state

import (
	"errors"
	"sync"
)

type entry struct {
	Name    string
	Version uint64
}

// StateImpl maintain metadata in memory
type StateImpl struct {
	sync.Mutex
	metadata map[string]entry
}

// NewState is a support function
func NewState() *StateImpl {
	return &StateImpl{
		metadata: map[string]entry{},
	}
}

// ChangeStatus update the state in the metadata, which will check the state machine is right.
func (state *StateImpl) ChangeStatus(token string, name string) (string, string, uint64, uint64, error) {
	state.Lock()
	defer state.Unlock()
	old, exist := state.metadata[token]
	if exist {
		state.metadata[token] = entry{Name: name, Version: old.Version + 1}
		return name, old.Name, old.Version + 1, old.Version, nil
	}
	state.metadata[token] = entry{Name: name, Version: 1}
	return name, "", 1, 0, nil
}

// Query returns the host
func (state *StateImpl) Query(token string) (string, uint64, error) {
	value, exist := state.metadata[token]
	if exist {
		return value.Name, value.Version, nil
	}
	return "", 0, errors.New("token area not exist")
}
