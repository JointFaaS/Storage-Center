package state

import (
	"errors"
)

// HostImpl maintain all hosts mapping name <=> host
type HostImpl struct {
	hosts map[string]string
}

// Insert a new host(for register)
func (h *HostImpl) Insert(host string, name string) error{
	_, exist := h.hosts[name]
	if (exist) {
		return errors.New("The name already exists")
	}
	h.hosts[name] = host
	return nil
}

// Query use name to get host
func (h *HostImpl) Query(name string) (string, error) {
	value, exist := h.hosts[name]
	if (exist) {
		return value, nil
	}
	return "", errors.New("The name is not found")
}

// Delete the name in the hosts
func (h *HostImpl) Delete(name string) error {
	_, exist := h.hosts[name]
	if (!exist) {
		return errors.New("The name is not found");
	}
	delete(h.hosts, name)
	return nil
}