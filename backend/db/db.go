package db

import "sync"

type store struct {
	mu    sync.RWMutex
	users map[string][]byte
}

func NewStore() *store {
	return &store{users: make(map[string][]byte)}
}

func (s *store) Add(email string, hash []byte) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[email]; exists {
		return false
	}
	s.users[email] = hash
	return true
}

func (s *store) Get(email string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.users[email]
	return h, ok
}
