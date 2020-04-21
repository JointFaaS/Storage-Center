package client


import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
)

// ClientStateImpl maintains state which stores keys you hold
type ClientStateImpl struct {
	holds mapset.Set
}

// Delete by token
func (s *ClientStateImpl) Delete(token string) {
	s.holds.Remove(token)
}

// Query by token
func (s *ClientStateImpl) Query(token string) bool {
	fmt.Printf("query state value: %v\n", s.holds)
	return s.holds.Contains(token)
}

// Add a token
func (s *ClientStateImpl) Add(token string) bool {
	result := s.holds.Add(token)
	fmt.Printf("add state value: %v\n", s.holds)
	return result
}