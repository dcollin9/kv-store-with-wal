package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Initialize initializes the in-memory key-value store by recovering data from the WAL and initializes the wal
// This function should be called once at application startup
func Initialize() error {
	// Log that we're initializing
	fmt.Println("Initializing key-value store...")

	// Clear the in-memory store first (in case of reinitialization)
	for k := range KVStore {
		delete(KVStore, k)
	}

	// Recover data from the WAL
	if err := RecoverFromWAL(); err != nil {
		return fmt.Errorf("failed to recover from WAL: %w", err)
	}

	// Initialize the global WAL file for writing
	if err := InitWAL(); err != nil {
		return fmt.Errorf("failed to initialize WAL: %w", err)
	}

	// Log how many items were recovered
	fmt.Printf("Recovery complete. Loaded %d items from WAL.\n", len(KVStore))
	fmt.Printf("KVstore: %+v", KVStore)

	return nil
}

// RecoverFromWAL reads the WAL file and populates the in-memory KVStore
// It reads the file in chunks of 100 bytes, handling newline-delimited values
func RecoverFromWAL() error {
	// Open the WAL file for reading
	walPath := filepath.Join(WALFileName)
	wal, err := os.OpenFile(walPath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			// No WAL file exists yet, that's okay
			return nil
		}
		return fmt.Errorf("failed to open WAL file for recovery: %w", err)
	}
	defer wal.Close()

	// Use a buffer to accumulate partial records
	var buffer bytes.Buffer

	// Read in chunks of 100 bytes
	// Note - if we have non ASCII characters (1 byte each), we could potentially have issues here,
	// esp if a special character, e.g., "é", that is more than a single byte. If our 100 byte chunk only includes one of the bytes,
	// then we'll have an invalid utf-8

	chunk := make([]byte, 100)
	for {
		bytesRead, err := wal.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			return fmt.Errorf("error reading from WAL: %w", err)
		}

		// Append the chunk to our buffer
		buffer.Write(chunk[:bytesRead])

		// Process complete records from the buffer
		processBuffer(&buffer)
	}

	// Process any remaining data in the buffer
	if buffer.Len() > 0 {
		processBuffer(&buffer)
	}

	return nil
}

// processBuffer extracts complete records from the buffer and updates the KVStore
func processBuffer(buffer *bytes.Buffer) {
	data := buffer.String()

	// Find the last complete record (ending with a newline)
	lastNewlineIndex := strings.LastIndex(data, "\n")
	if lastNewlineIndex == -1 {
		// No complete record in buffer yet
		return
	}

	// Extract the complete portion
	completeData := data[:lastNewlineIndex+1]

	// Process the complete data line by line
	scanner := bufio.NewScanner(strings.NewReader(completeData))
	for scanner.Scan() {
		record := scanner.Text()
		if record == "" {
			continue
		}

		parts := strings.SplitN(record, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			KVStore[key] = value
		}
	}

	// Keep the incomplete portion in the buffer, write the rest
	if lastNewlineIndex < len(data)-1 {
		buffer.Reset()
		buffer.WriteString(data[lastNewlineIndex+1:])
	} else {
		buffer.Reset()
	}
}
