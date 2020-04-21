package state

import (
	"errors"
	"fmt"
)

// StatusLine store host and chan
type StatusLine struct {
	host           string
	invalidChannel chan string
}

// HostImpl maintain all hosts mapping name <=> host
type HostImpl struct {
	hosts map[string]StatusLine
}

// NewHost is an support function
func NewHost() *HostImpl {
	return &HostImpl{
		hosts: make(map[string]StatusLine),
	}
}

// Insert a new host(for register)
func (h *HostImpl) Insert(host string, name string) error {
	_, exist := h.hosts[name]
	if exist {
		return errors.New("The name already exists")
	}
	fmt.Printf("Insert %v\n", h.hosts)
	fmt.Printf("Host %v, Name %v\n", host, name)
	h.hosts[name] = StatusLine{
		host:           host,
		invalidChannel: make(chan string, 100),
	}
	return nil
}

// Query use name to get host
func (h *HostImpl) Query(name string) (string, error) {
	value, exist := h.hosts[name]
	if exist {
		return value.host, nil
	}
	return "", errors.New("The name is not found")
}

// Delete the name in the hosts
func (h *HostImpl) Delete(name string) error {
	_, exist := h.hosts[name]
	if !exist {
		return errors.New("The name is not found")
	}
	delete(h.hosts, name)
	return nil
}

// GetChan return channel for invalid communication
func (h *HostImpl) GetChan(name string) (chan string, error) {
	value, exist := h.hosts[name]
	if exist {
		return value.invalidChannel, nil
	}
	return nil, errors.New("The name is not found")
}
