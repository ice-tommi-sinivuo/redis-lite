package storage

import (
	"sync"
	"testing"
)

func TestMemoryStore_Set(t *testing.T) {
	store := NewMemoryStore()

	err := store.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set() returned error: %v", err)
	}

	// Verify the value was stored
	value, exists := store.Get("key1")
	if !exists {
		t.Error("Key was not stored")
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%s'", value)
	}
}

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemoryStore()

	// Test getting non-existent key
	value, exists := store.Get("nonexistent")
	if exists {
		t.Error("Expected key to not exist")
	}
	if value != "" {
		t.Errorf("Expected empty value for non-existent key, got '%s'", value)
	}

	// Test getting existing key
	store.Set("key1", "value1")
	value, exists = store.Get("key1")
	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%s'", value)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()

	// Test deleting non-existent key
	deleted := store.Delete("nonexistent")
	if deleted {
		t.Error("Expected delete to return false for non-existent key")
	}

	// Test deleting existing key
	store.Set("key1", "value1")
	deleted = store.Delete("key1")
	if !deleted {
		t.Error("Expected delete to return true for existing key")
	}

	// Verify key is gone
	_, exists := store.Get("key1")
	if exists {
		t.Error("Key should not exist after deletion")
	}
}

func TestMemoryStore_Exists(t *testing.T) {
	store := NewMemoryStore()

	// Test non-existent key
	if store.Exists("nonexistent") {
		t.Error("Expected key to not exist")
	}

	// Test existing key
	store.Set("key1", "value1")
	if !store.Exists("key1") {
		t.Error("Expected key to exist")
	}

	// Test after deletion
	store.Delete("key1")
	if store.Exists("key1") {
		t.Error("Expected key to not exist after deletion")
	}
}

func TestMemoryStore_Size(t *testing.T) {
	store := NewMemoryStore()

	// Test empty store
	if store.Size() != 0 {
		t.Errorf("Expected size 0, got %d", store.Size())
	}

	// Test after adding keys
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	if store.Size() != 2 {
		t.Errorf("Expected size 2, got %d", store.Size())
	}

	// Test after deletion
	store.Delete("key1")
	if store.Size() != 1 {
		t.Errorf("Expected size 1, got %d", store.Size())
	}
}

func TestMemoryStore_Clear(t *testing.T) {
	store := NewMemoryStore()

	// Add some keys
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Set("key3", "value3")

	// Clear the store
	store.Clear()

	// Verify all keys are gone
	if store.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", store.Size())
	}

	if store.Exists("key1") || store.Exists("key2") || store.Exists("key3") {
		t.Error("Keys should not exist after clear")
	}
}

func TestMemoryStore_OverwriteValue(t *testing.T) {
	store := NewMemoryStore()

	// Set initial value
	store.Set("key1", "value1")

	// Overwrite with new value
	store.Set("key1", "value2")

	// Verify new value
	value, exists := store.Get("key1")
	if !exists {
		t.Error("Key should exist")
	}
	if value != "value2" {
		t.Errorf("Expected value 'value2', got '%s'", value)
	}

	// Verify size is still 1
	if store.Size() != 1 {
		t.Errorf("Expected size 1, got %d", store.Size())
	}
}

func TestMemoryStore_EmptyValues(t *testing.T) {
	store := NewMemoryStore()

	// Test empty string value
	store.Set("empty", "")
	value, exists := store.Get("empty")
	if !exists {
		t.Error("Key with empty value should exist")
	}
	if value != "" {
		t.Errorf("Expected empty string, got '%s'", value)
	}
}

func TestMemoryStore_ThreadSafety(t *testing.T) {
	store := NewMemoryStore()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				key := "key"
				value := "value"

				// Perform various operations
				store.Set(key, value)
				store.Get(key)
				store.Exists(key)
				store.Size()

				if j%10 == 0 {
					store.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()

	// The test passes if no race conditions are detected
	// (run with -race flag to verify)
}
