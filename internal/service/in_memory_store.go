// Package service provides business logic for key-value store operations.
package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	KVStore[key] = val

	// now write the key/val pair to the Write Ahead Log (WAL)
	walPath := filepath.Join(WALFileName)
	wal, err := os.OpenFile(walPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer wal.Close()

	// Format with comma delimiter to match recovery function
	data := []byte(fmt.Sprintf("%s=%s,", key, val))
	if _, err := wal.Write(data); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
