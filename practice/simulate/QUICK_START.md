# DNA Search Algorithms - Quick Start

## 🚀 Run Any Demo in 3 Seconds

All demos are now in separate folders with `package main`!

### 1. N-gram Search (Fast repeated queries)
```bash
cd dna_ngram_search_demo && go run dna_ngram_search.go
```

### 2. Seed-and-Extend (BLAST-style, production-ready)
```bash
cd dna_seed_extend_demo && go run dna_seed_extend.go
```

### 3. Performance Benchmark (See the speedup!)
```bash
cd dna_seed_extend_demo && go run seed_extend_benchmark.go
```

### 4. Strategy Comparison (Learn all 5 approaches)
```bash
cd dna_wildcard_strategies_demo && go run wildcard_comparison_demo.go
```

### 5. Original Algorithms (Educational)
```bash
cd dna_string_search && go run dna_string_search.go
```

## 📊 Example Output

### Seed-and-Extend Demo
```
Query: ATTR (R = A or G)
Extracted seed: 'ATT' at offset 0

Seed 'ATT' found at:
  - Sequence 'seq1' position 1 → ATTA ✓
  - Sequence 'seq1' position 6 → ATTA ✓
  - Sequence 'seq2' position 3 → ATTA ✓
  - Sequence 'seq3' position 4 → ATTG ✓
```

### Benchmark Results
```
Query "ATTR" on 1MB database:
  Seed-and-Extend:  2.7ms
  Naive:           37.0ms
  Speedup:         13.7x faster! ⚡
```

## 🎯 Which One Should I Run?

| Goal | Run This |
|------|----------|
| Learn basics | `dna_string_search/` |
| See wildcard expansion | `dna_ngram_search_demo/` |
| **Production algorithm** | **`dna_seed_extend_demo/`** ⭐ |
| Compare approaches | `dna_wildcard_strategies_demo/` |
| Measure performance | `seed_extend_benchmark.go` |

## 📚 More Info

- [README.md](README.md) - Full overview
- [README_FOLDER_STRUCTURE.md](dna_ngram_search/README_FOLDER_STRUCTURE.md) - Detailed guide
- Each folder has its own documentation!

---

**All ready to run! No dependencies needed.** 🎉
