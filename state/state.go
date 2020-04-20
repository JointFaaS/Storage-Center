package state

import (
	"errors"
)

// StateImpl maintain metadata in memory
type StateImpl struct {
	metadata map[string]string
}

// NewState is a support function
func NewState() *StateImpl {
	return &StateImpl{
		metadata: map[string]string{},
	}
}

// ChangeStatus update the state in the metadata, which will check the state machine is right.
func (state *StateImpl) ChangeStatus(token string, name string) (string, string, error) {
	old, exist := state.metadata[token]
	state.metadata[token] = name
	if exist {
		return name, old, nil
	}
	return name, "", nil

}

// Query returns the host
func (state *StateImpl) Query(token string) (string, error) {
	value, exist := state.metadata[token]
	if exist {
		return value, nil
	}
	return "", errors.New("token area not exist")
}
