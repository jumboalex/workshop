# Mini KV Store

A thread-safe key-value store implementation in Go with snapshot capabilities and an interactive CLI.

## Project Structure

```
kv_store/
├── store/
│   ├── kv_store.go          # Core KV store implementation
│   └── kv_store_test.go     # Comprehensive unit tests (88.2% coverage)
├── cmd/
│   ├── cli/
│   │   └── main.go          # Interactive CLI
│   └── demo/
│       └── main.go          # Automated demo/integration tests
├── bin/                     # Compiled binaries
│   ├── kv-cli
│   └── kv-demo
├── go.mod                   # Go module definition
└── README.md                # Documentation (this file)
```

## Operations

- `Put(key, value)` - Store a key-value pair
- `Get(key)` - Retrieve value for a key
- `CountValue(value)` - Count how many keys have the given value
- `Checkpoint()` - Create a snapshot of current state
- `Revert()` - Restore to last checkpoint
- `SaveToDisk(filename)` - Persist state to disk (JSON)
- `LoadFromDisk(filename)` - Load state from disk
- `GetCheckpointCount()` - Get number of checkpoints
- `GetAllData()` - Get copy of all data and value counts

## Quick Start

```bash
# Build both programs
go build -o bin/kv-demo ./cmd/demo
go build -o bin/kv-cli ./cmd/cli

# Run the demo tests
./bin/kv-demo

# Run the interactive CLI
./bin/kv-cli

# Or run directly without building
go run ./cmd/demo
go run ./cmd/cli
```

## Testing

The project includes comprehensive unit tests with 88.2% code coverage.

```bash
# Run all tests
go test ./store/

# Run tests with verbose output
go test -v ./store/

# Run tests with race detection
go test -race ./store/

# Run tests with coverage report
go test -cover ./store/

# Generate detailed coverage report
go test -coverprofile=coverage.out ./store/
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./store/
```

### Test Coverage

- **TestPutAndGet** - Basic put and get operations
- **TestPutOverwrite** - Overwriting existing keys
- **TestCountValue** - Value counting functionality
- **TestCheckpointAndRevert** - Snapshot and restore
- **TestMultipleCheckpoints** - Multiple checkpoint levels
- **TestRevertWithoutCheckpoint** - Error handling
- **TestSaveAndLoad** - Disk persistence
- **TestGetAllData** - Data retrieval methods
- **TestConcurrentPut** - Concurrent write operations
- **TestConcurrentGetAndPut** - Concurrent read/write mix
- **TestConcurrentCheckpoint** - Concurrent snapshots
- **TestEmptyStore** - Empty store edge cases
- **TestCheckpointPreservesValueCounts** - Value count preservation
- **TestLoadInvalidFile** - Error handling for invalid files

### Benchmark Results

```
BenchmarkPut-12            2142691    686.3 ns/op    309 B/op    4 allocs/op
BenchmarkGet-12           16919719     71.1 ns/op     13 B/op    1 allocs/op
BenchmarkCountValue-12    77416527     15.1 ns/op      0 B/op    0 allocs/op
BenchmarkCheckpoint-12      179492   6159.0 ns/op   8552 B/op    9 allocs/op
BenchmarkRevert-12          111004  11172.0 ns/op  17007 B/op   16 allocs/op
```

Performance highlights:
- **CountValue**: ~15 ns/op (O(1) with valueCount map)
- **Get**: ~71 ns/op (very fast lookups)
- **Put**: ~686 ns/op (includes valueCount bookkeeping)
- **No race conditions** detected in concurrent tests

## Implementation Details

### Snapshot Structure

Uses a `Snapshot` struct to capture both data and valueCount at checkpoint time:

```go
type Snapshot struct {
    Data       map[string]string
    ValueCount map[string]int
}
```

**Benefits:**
- No need to rebuild valueCount after revert
- Complete state preservation
- Faster revert operations
- Consistent snapshots

### Current Implementation (With valueCount map)

This implementation maintains a separate `valueCount` map for O(1) value counting.

**Pros:**
- `CountValue()` is **O(1)** - instant lookup
- Much faster for frequent counting operations
- Excellent for high-performance scenarios
- No rebuild needed on `Revert` (uses stored valueCount from Snapshot)

**Cons:**
- More memory usage (extra map + snapshots store both maps)
- Must maintain consistency between data and valueCount during `Put`
- Slightly more complex `Put` operation

**Time Complexity:**
- `Put`: O(1) with valueCount bookkeeping
- `Get`: O(1)
- `CountValue`: **O(1)** ✨
- `Checkpoint`: O(n) - deep copies both maps
- `Revert`: O(n) - restores both maps (no rebuild needed!)

### Alternative: Without valueCount map

**Pros:**
- Simpler data structure
- Less memory overhead
- No synchronization needed between maps

**Cons:**
- `CountValue()` would be **O(n)** - must scan entire map
- Slower for frequent value counting operations

**When to use:**
- Use current implementation (with valueCount) for production use
- Use simple implementation (without valueCount) only if `CountValue` is rarely called

## Thread Safety

The KV store uses `sync.RWMutex` for thread-safe concurrent access:

- **Read operations** (`Get`, `CountValue`, `SaveToDisk`, `Print`) use `RLock()` - multiple readers can access simultaneously
- **Write operations** (`Put`, `Checkpoint`, `Revert`, `LoadFromDisk`) use `Lock()` - exclusive access

**Concurrency features:**
- Safe for concurrent reads and writes from multiple goroutines
- No race conditions (verified with `go run -race`)
- Efficient reader/writer lock pattern
- All operations are atomic

## Interactive CLI

The package includes a command-line interface for easy interaction with the KV store.

### Building and Running

```bash
# Build the CLI
go build -o bin/kv-cli ./cmd/cli

# Run the CLI
./bin/kv-cli

# Or run directly
go run ./cmd/cli
```

### Available Commands

| Command | Aliases | Description | Example |
|---------|---------|-------------|---------|
| `put <key> <value>` | | Store a key-value pair | `put name Alice` |
| `get <key>` | | Retrieve value for key | `get name` |
| `count <value>` | | Count keys with value | `count active` |
| `checkpoint` | `cp` | Create snapshot | `checkpoint` |
| `revert` | `rv` | Revert to last checkpoint | `revert` |
| `save [file]` | | Save to disk | `save store.json` |
| `load [file]` | | Load from disk | `load store.json` |
| `list` | `ls` | Show all data | `list` |
| `clear` | | Clear the store | `clear` |
| `help` | `?` | Show help | `help` |
| `exit` | `quit`, `q` | Exit CLI | `exit` |

### Flag-Based CLI Example

```bash
# Put values
$ kv-cli put --key name --value Alice
✅ Set 'name' = 'Alice'

$ kv-cli put --key message --value "hello world"
✅ Set 'message' = 'hello world'

# Get values
$ kv-cli get --key name
✅ 'name' = 'Alice'

# Count values
$ kv-cli count --value Alice
✅ Value 'Alice' appears 1 time(s)

# Checkpoint and revert
$ kv-cli checkpoint
✅ Checkpoint created (total: 1)

$ kv-cli put --key name --value Bob
✅ Set 'name' = 'Bob'

$ kv-cli revert
✅ Reverted to last checkpoint

# List all data
$ kv-cli list
Current Key-Value Pairs:
------------------------
  name = Alice
  message = hello world

Value Counts:
-------------
  'Alice' → 1
  'hello world' → 1

Total keys: 2
Checkpoints: 0
```

**Auto-Persistence:** The CLI automatically saves data to `.kv_store.json` after each operation and loads it on startup, so data persists between invocations.

**Script-Friendly:** Perfect for shell scripts and automation:
```bash
#!/bin/bash
kv-cli put --key deployment_time --value "$(date)"
kv-cli put --key version --value "v1.2.3"
kv-cli checkpoint
```

## Configuration Decisions

**Do you need locks?**
- **YES** if multiple goroutines will access the store
- **NO** if single-threaded or externally synchronized
- Without locks: simpler code but not thread-safe
- With locks: thread-safe but slight performance overhead

## Usage Example

```go
kv := NewKVStore()

// Basic operations
kv.Put("user1", "active")
kv.Put("user2", "active")
kv.Put("user3", "inactive")

value, ok := kv.Get("user1")  // "active", true
count := kv.CountValue("active")  // 2

// Snapshots
kv.Checkpoint()
kv.Put("user1", "inactive")
kv.Revert()  // user1 is back to "active"

// Persistence
kv.SaveToDisk("store.json")
kv.LoadFromDisk("store.json")

// Concurrent access (thread-safe)
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        kv.Put(fmt.Sprintf("key%d", id), "value")
        kv.Get(fmt.Sprintf("key%d", id))
        kv.CountValue("value")
    }(i)
}
wg.Wait()
```
