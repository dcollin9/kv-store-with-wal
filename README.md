# kv-store-with-wal
This repo holds a WAL-backed KV store, implemented in Go.

## Purpose
Build a basic in-memory key-value store (in Go), with a write-ahead log on disk to ensure durability and crash recovery.

## Functionality
Accept SET key value and GET key commands.

On SET, write the operation to the WAL (append-only file) before applying it to the in-memory map.

On startup, replay the WAL to reconstruct the in-memory state.

Keep it single-threaded — no need for concurrency yet.

This mimics the basic architecture of Postgres or InnoDB, where a WAL guarantees durability even if the in-memory state is lost.

## How to Run

Run `go run cmd/api/main.go` to start the application.

To test:

Run `go run scripts/get/main.go ex_key` to get a value from `ex_key`

Run `go run scripts/set/main.go ex_key ex_value` to set the example key value pair in memory.


## Linters
Run `golangci-lint help linters` to check for all avialable linters.

Run the linting using: `golangci-lint run`

## Application Architecture

This application implements a crash-safe, in-memory key-value store with the following architecture:

### Core Components

**In-Memory Store**: A simple Go map (`KVStore`) that holds all key-value pairs in memory for fast access.

**Write-Ahead Log (WAL)**: A persistent text file (`wal.txt`) that records all write operations before they're applied to the in-memory store, ensuring durability and crash recovery.

**HTTP API**: RESTful endpoints for reading and writing data:
- `GET /v1/{key}` - Retrieve a value by key
- `POST /v1/write` - Store a key-value pair

### WAL Implementation Details

**Startup Recovery**: On application startup, the WAL file is read in 100-byte chunks to reconstruct the in-memory store from persisted operations. Records are newline-delimited in the format `key=value`.

**Global WAL Connection**: The WAL file is opened once at startup and kept open throughout the application lifecycle, eliminating the overhead of repeated file open/close operations.

**Concurrent Write Safety**: All WAL writes are protected by a mutex to prevent race conditions when multiple goroutines attempt to write simultaneously.

**Periodic Syncing**: A background goroutine uses a ticker to flush the WAL to disk every 5 seconds, balancing performance with durability guarantees.

**Write-First Semantics**: All SET operations write to the WAL first, then update the in-memory store only after successful WAL persistence, ensuring consistency.

**Graceful Shutdown**: On receiving shutdown signals (SIGINT/SIGTERM), the application performs a final WAL sync and properly closes the file before terminating.

### Performance Characteristics

- **Reads**: O(1) from in-memory map
- **Writes**: O(1) to memory + append to WAL file
- **Startup**: O(n) where n is the number of records in the WAL
- **Durability**: Data is guaranteed durable within 5 seconds (configurable sync interval)

This design provides a good balance of performance, simplicity, and crash safety suitable for applications requiring persistent key-value storage with fast access patterns.