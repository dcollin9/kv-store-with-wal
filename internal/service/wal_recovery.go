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

// InitStore initializes the in-memory key-value store by recovering data from the WAL
// This function should be called once at application startup
func InitStore() error {
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

	// Log how many items were recovered
	fmt.Printf("Recovery complete. Loaded %d items from WAL.\n", len(KVStore))
	fmt.Printf("KVstore: %+v", KVStore)

	return nil
}

// RecoverFromWAL reads the WAL file and populates the in-memory KVStore
// It reads the file in chunks of 100 bytes, handling comma-delimited values
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

	// Find the last complete record (ending with a comma)
	lastCommaIndex := strings.LastIndex(data, ",")
	if lastCommaIndex == -1 {
		// No complete record in buffer yet
		return
	}

	// Extract the complete portion
	completeData := data[:lastCommaIndex+1]

	// Create a scanner to process each complete record
	scanner := bufio.NewScanner(strings.NewReader(completeData))
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Look for comma
		if i := bytes.IndexByte(data, ','); i >= 0 {
			// We have a complete record
			return i + 1, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated record
		if atEOF {
			return len(data), data, nil
		}

		// Request more data
		return 0, nil, nil
	})

	// Process each record and update KVStore
	for scanner.Scan() {
		record := scanner.Text()
		parts := strings.SplitN(record, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			KVStore[key] = value
		}
	}

	// Keep only the incomplete portion in the buffer
	if lastCommaIndex < len(data)-1 {
		buffer.Reset()
		buffer.WriteString(data[lastCommaIndex+1:])
	} else {
		buffer.Reset()
	}
}
