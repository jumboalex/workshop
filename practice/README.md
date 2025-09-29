# Practice Problems

A collection of algorithm and data structure problems organized in Go packages.

## Project Structure

```
practice/
├── algorithms/
│   ├── array/           # Array-based problems
│   ├── string/          # String manipulation problems
│   ├── linkedlist/      # Linked list problems
│   └── tree/            # Binary tree problems
├── cmd/
│   └── runner/          # CLI to run examples
├── Makefile             # Build automation
└── README.md
```

## Quick Start

### Prerequisites
- Go 1.16 or higher

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage report
make coverage

# Run specific package tests
go test -v ./algorithms/array
go test -v ./algorithms/string
go test -v ./algorithms/linkedlist
go test -v ./algorithms/tree
```

### Building
```bash
# Build the runner executable
make build

# Build and run examples
make run
```

### Other Commands
```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Check code (format, vet, test)
make check

# Show all available commands
make help
```

## Packages

### Array Problems (`algorithms/array`)
- **PlusOne**: Add one to number represented as array of digits
- **RemoveDuplicates**: Remove duplicates from sorted array in-place
- **CanPlaceFlowers**: Determine if flowers can be placed without adjacency
- **KidsWithCandies**: Find kids who can have max candies with extra
- **MaxOperations**: Find max operations to form pairs summing to k
- **LongestOnes**: Longest subarray of 1s after flipping k zeros
- **LongestSubarray**: Longest subarray of 1s after deleting one element

### String Problems (`algorithms/string`)
- **GcdOfStrings**: Find greatest common divisor string
- **MergeAlternately**: Merge two strings alternately
- **AddBinary**: Add two binary strings
- **MultiplyString**: Multiply two numbers as strings
- **MaxVowels**: Maximum vowels in substring of length k
- **LengthOfLongestSubstringKDistinct**: Longest substring with k distinct chars

### Linked List Problems (`algorithms/linkedlist`)
- **ReorderList**: Reorder list as L0→Ln→L1→Ln-1→L2→Ln-2...
- **CopyRandomList**: Deep copy list with random pointers

### Tree Problems (`algorithms/tree`)
- **Flatten**: Flatten binary tree to linked list in-place
- **IsValidBST**: Validate if tree is binary search tree

## Usage Examples

### As a Library
```go
import (
    arrayproblems "github.com/jumbo/workshop/practice/algorithms/array"
    stringproblems "github.com/jumbo/workshop/practice/algorithms/string"
)

func main() {
    // Array problem
    result := arrayproblems.PlusOne([]int{1, 2, 9})
    fmt.Println(result) // [1 3 0]

    // String problem
    merged := stringproblems.MergeAlternately("abc", "pqr")
    fmt.Println(merged) // "apbqcr"
}
```

### Running the CLI
```bash
# Build and run all examples
make run

# Or directly
go run cmd/runner/main.go

# Show help
./bin/runner --help
```

## Testing

All functions have comprehensive unit tests with multiple test cases covering:
- Normal cases
- Edge cases
- Empty inputs
- Large inputs

Test coverage can be generated with:
```bash
make coverage
open coverage.html  # View in browser
```

## Development Workflow

1. **Add a new problem**:
   - Add function to appropriate package in `algorithms/`
   - Write tests in corresponding `*_test.go` file
   - Run `make test` to verify

2. **Format and check**:
   ```bash
   make check  # Formats, vets, and tests
   ```

3. **Add to runner**:
   - Update `cmd/runner/main.go` to showcase your problem
   - Run `make run` to test

## Makefile Targets

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make coverage` | Generate coverage report |
| `make build` | Build runner executable |
| `make run` | Build and run application |
| `make fmt` | Format all code |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint |
| `make clean` | Remove build artifacts |
| `make bench` | Run benchmarks |
| `make tidy` | Tidy dependencies |
| `make check` | Format, vet, and test |
| `make help` | Show help |

## License

This is a personal practice repository.