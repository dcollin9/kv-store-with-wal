package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// Global WAL file handle
	walFile  *os.File
	walMutex sync.Mutex

	// Channel to signal shutdown
	shutdownChan = make(chan struct{})
	syncTicker   *time.Ticker
)

// InitWAL opens the WAL file globally and starts the periodic sync routine
func InitWAL() error {
	walPath := filepath.Join(WALFileName)
	var err error
	walFile, err = os.OpenFile(walPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open WAL file: %w", err)
	}

	// Start periodic sync routine (every 5 seconds)
	// Ticker vs for loop with sleep:
	// Ticker executes no matter how long the previous sync takes, automatically handles timing resources, is select statement friendly. But, slightly more memory overhead
	// For loop is simpler, has adaptive timing, uses less memory. But is has drift, is harder to cancel, and is blocking.
	// The select is better for shutdown, has bounded memory, and is the idiomatic approach for periodic tasks
	syncTicker = time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-syncTicker.C:
				fmt.Println("syncing WAL...")
				walMutex.Lock()
				if walFile != nil {
					walFile.Sync()
				}
				walMutex.Unlock()
			case <-shutdownChan:
				return
			}
		}
	}()

	return nil
}

// CloseWAL performs final sync and closes the WAL file
func CloseWAL() error {
	// Stop the sync ticker
	if syncTicker != nil {
		syncTicker.Stop()
	}

	// Signal shutdown to the sync goroutine
	close(shutdownChan)

	walMutex.Lock()
	defer walMutex.Unlock()

	if walFile != nil {
		// Final sync before closing
		if err := walFile.Sync(); err != nil {
			return fmt.Errorf("failed to sync WAL on close: %w", err)
		}

		if err := walFile.Close(); err != nil {
			return fmt.Errorf("failed to close WAL file: %w", err)
		}
		walFile = nil
	}

	return nil
}
