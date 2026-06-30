package store

import (
	"sync"
)

// MemStore is a thread-safe, in-memory implementation of the Store interface.
type MemStore struct {
	mu    sync.RWMutex
	items map[string]string
}

// NewMemStore initializes a new MemStore.
func NewMemStore() *MemStore {
	return &MemStore{
		items: make(map[string]string),
	}
}

func (m *MemStore) Set(key string, value string) error {
	m.mu.Lock()         
	defer m.mu.Unlock() 

	m.items[key] = value
	return nil
}

func (m *MemStore) Get(key string) (string, error) {
	m.mu.RLock()         
	defer m.mu.RUnlock()

	value, exists := m.items[key]
	if !exists {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (m *MemStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}