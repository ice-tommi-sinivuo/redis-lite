package storage

import (
	"sync"
)

// Store defines the interface for data storage operations
type Store interface {
	// Set stores a key-value pair
	Set(key, value string) error

	// Get retrieves a value by key
	Get(key string) (string, bool)

	// Delete removes a key-value pair
	Delete(key string) bool

	// Exists checks if a key exists
	Exists(key string) bool

	// Size returns the number of stored keys
	Size() int

	// Clear removes all keys
	Clear()
}

// MemoryStore implements Store interface with in-memory storage
type MemoryStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair
func (s *MemoryStore) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
	return nil
}

// Get retrieves a value by key, returns value and whether the key exists
func (s *MemoryStore) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Delete removes a key-value pair, returns true if key existed
func (s *MemoryStore) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
	}
	return exists
}

// Exists checks if a key exists
func (s *MemoryStore) Exists(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.data[key]
	return exists
}

// Size returns the number of stored keys
func (s *MemoryStore) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.data)
}

// Clear removes all keys
func (s *MemoryStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = make(map[string]string)
}
