package client

import (
	"fmt"
	"sync"

	"github.com/JointFaaS/Storage-Center/enum"
)

type entry struct {
	Version uint64
	Status  enum.Policy
}

// ClientStateImpl maintains state which stores keys you hold
type ClientStateImpl struct {
	sync.Mutex
	holds map[string]entry
}

// NewClientStateImpl is a support function
func NewClientStateImpl() *ClientStateImpl {
	return &ClientStateImpl{
		holds: make(map[string]entry),
	}
}

// Delete by token
func (s *ClientStateImpl) Delete(token string) {
	s.Lock()
	defer s.Unlock()
	delete(s.holds, token)
}

// Query by token
func (s *ClientStateImpl) Query(token string) bool {
	s.Lock()
	defer s.Unlock()
	fmt.Printf("query state value: %v, token %v\n", s.holds, token)
	_, exist := s.holds[token]
	return exist
}

// GetVersion by token
func (s *ClientStateImpl) GetVersion(token string) uint64 {
	s.Lock()
	defer s.Unlock()
	fmt.Printf("getversion state value: %v, token %v\n", s.holds, token)
	value, exist := s.holds[token]
	if exist {
		return value.Version
	}
	return 0
}

// GetStatus by token return Policy
func (s *ClientStateImpl) GetStatus(token string) enum.Policy {
	s.Lock()
	defer s.Unlock()
	fmt.Printf("getversion state value: %v, token %v\n", s.holds, token)
	value, exist := s.holds[token]
	if exist {
		return value.Status
	}
	return enum.PolicyInvalid
}

// Add a token
func (s *ClientStateImpl) Add(token string, version uint64, status enum.Policy) bool {
	s.Lock()
	defer s.Unlock()
	s.holds[token] = entry{Version: version, Status: status}
	fmt.Printf("add state value: %v\n", s.holds)
	return true
}
