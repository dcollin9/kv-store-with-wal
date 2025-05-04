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