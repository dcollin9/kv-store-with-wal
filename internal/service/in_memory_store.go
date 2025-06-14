// Package service provides business logic for key-value store operations.
package service

import (
	"context"
	"errors"
	"fmt"
)

// WALFileName is the name of the write-ahead log file
const WALFileName = "wal.txt"

// Standard errors for the service.
var (
	// ErrNotFound is returned when a requested key doesn't exist.
	ErrNotFound = errors.New("key not found")

	// KVStore holds all key-value pairs in memory.
	KVStore = map[string]string{}
)

// KVPair represents a key-value pair for JSON serialization.
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Get retrieves a value from the store by its key.
// It returns the value if found, or an error if the key doesn't exist.
func Get(_ context.Context, key string) (string, error) {
	if v, ok := KVStore[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("get operation failed: %w", ErrNotFound)
}

// Set stores a key-value pair in the in-memory store.
// It overwrites any existing value for the same key.
func Set(_ context.Context, key, val string) error {
	// Write to WAL first (before updating in-memory store)
	walMutex.Lock()

	data := []byte(fmt.Sprintf("%s=%s\n", key, val))
	if _, err := walFile.Write(data); err != nil {
		walMutex.Unlock()
		return fmt.Errorf("failed to write to WAL: %w", err)
	}
	walMutex.Unlock()

	// Only update in-memory store after successful WAL write
	KVStore[key] = val

	return nil
}
