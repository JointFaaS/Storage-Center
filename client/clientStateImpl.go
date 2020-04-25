package client

import (
	"fmt"
	"sync"
)

// ClientStateImpl maintains state which stores keys you hold
type ClientStateImpl struct {
	sync.Mutex
	holds map[string]uint64
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
		return value
	}
	return 0
}

// Add a token
func (s *ClientStateImpl) Add(token string, version uint64) bool {
	s.Lock()
	defer s.Unlock()
	s.holds[token] = version
	fmt.Printf("add state value: %v\n", s.holds)
	return true
}
