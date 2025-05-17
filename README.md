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