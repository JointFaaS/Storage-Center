package client 

import (
	"sync"
	"errors"
)

type StorageImpl struct {
	sync.RWMutex
	storage map[string]string
}

func (s *StorageImpl)Set(token string, value string) {
	s.Lock()
	s.storage[token] = value
	s.Unlock()
}

func (s *StorageImpl)Get(token string) string {
	s.RLock()
	value := s.storage[token]
	s.RUnlock()
	return value
}

func (s *StorageImpl)Delete(token string) error {
	s.Lock()
	_, exist := s.storage[token]
	if !exist {
		return errors.New("the value of " + token + "is not exist")
	}
	delete(s.storage, token)
	s.Unlock()
	return nil
}