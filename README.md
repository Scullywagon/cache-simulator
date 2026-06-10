# cache-protocols

Small Go playground for modeling cache behavior against a byte-addressable memory
store. The repo currently focuses on a simple direct-mapped cache, address
decoding, and deterministic simulator setup for tests.

## Current State

- `cache`: cache lines, reads, and inserts with tag/valid checks.
- `memory`: byte-backed memory with bounded block reads.
- `system`: address layout calculation, address decoding, and cached reads that
  fetch from memory on misses.
- `simulator`: deterministic system construction with seeded memory data and
  cache warmup modes: cold, full, partial, and random.
- `main.go`: placeholder executable that only prints `hello world`.

This is still experimental. The `Stats` fields are defined but not updated yet,
`memory.Write` is a stub, and there is not currently a real CLI or simulation
runner.

## Requirements

- Go `1.25.4` as declared in `go.mod`

## Useful Commands

```sh
go test ./...
```

If the default Go build cache is not writable in your environment, point it at a
local temporary directory:

```sh
GOCACHE=/tmp/cache-protocols-go-build go test ./...
```

## Layout

```text
.
├── cache/      # Direct-mapped cache model and tests
├── memory/     # Byte-addressable memory model
├── simulator/  # Seeded system setup and warmup behavior
├── system/     # Cache/memory controller and address decoding
└── main.go     # Placeholder entry point
```
