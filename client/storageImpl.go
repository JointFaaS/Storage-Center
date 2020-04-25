package client

import (
	"errors"
	"sync"
)

// StorageImpl is an type of storage using simple map and Mutex
type StorageImpl struct {
	sync.Mutex
	storage map[string]string
}

// Set key
func (s *StorageImpl) Set(token string, value string) {
	s.Lock()
	s.storage[token] = value
	s.Unlock()
}

// Get value by key
func (s *StorageImpl) Get(token string) (string, error) {
	s.Lock()
	value, exist := s.storage[token]
	s.Unlock()
	if exist {
		return value, nil
	}
	return "", errors.New("Not Found")
}

// Delete will delete key in map
func (s *StorageImpl) Delete(token string) error {
	s.Lock()
	_, exist := s.storage[token]
	if !exist {
		return nil
		// return errors.New("the value of " + token + " is not exist")
	}
	delete(s.storage, token)
	s.Unlock()
	return nil
}
