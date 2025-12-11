# DNA Search Algorithms - Folder Structure

This directory contains multiple DNA wildcard search implementations. Each has been organized into its own folder so they can all use `package main` and be run independently.

## 📁 Folder Organization

```
practice/simulate/
├── dna_string_search/           # Original string search algorithms
│   ├── dna_string_search.go     # Naive, KMP, Sliding Window
│   └── dna_string_search.cpp    # C++ version
│
├── dna_ngram_search_demo/       # N-gram indexing with wildcards
│   ├── dna_ngram_search.go      # Main implementation
│   ├── README.md                # Full documentation
│   ├── BENEFITS.md              # Performance analysis
│   └── WILDCARD_FIX.md          # How wildcard expansion works
│
├── dna_seed_extend_demo/        # Seed-and-extend (BLAST-style)
│   ├── dna_seed_extend.go       # Main implementation
│   ├── seed_extend_benchmark.go # Performance benchmarks
│   └── SEED_EXTEND.md           # Algorithm documentation
│
└── dna_wildcard_strategies_demo/  # Comparison of 8 strategies
    ├── wildcard_comparison_demo.go  # Demo implementations
    └── WILDCARD_STRATEGIES.md       # Strategy comparison guide
```

## 🚀 Quick Start Guide

### 1. N-gram Search (Best for: Known database, fast queries)

```bash
cd practice/simulate/dna_ngram_search_demo
go run dna_ngram_search.go
```

**Use when:**
- Searching a fixed database repeatedly
- Need O(1) seed lookups
- Wildcards are moderate (1-3 per query)

**Features:**
- Inverted index for fast lookups
- Wildcard expansion during query
- Handles multiple sequences

### 2. Seed-and-Extend (Best for: Large databases, many wildcards)

```bash
cd practice/simulate/dna_seed_extend_demo
go run dna_seed_extend.go              # Demo
go run seed_extend_benchmark.go         # Benchmarks
```

**Use when:**
- Database is very large (genomic scale)
- Queries have many wildcards (3+)
- Want proven BLAST-style algorithm

**Features:**
- Only checks promising positions
- Handles unlimited wildcards
- Used by BLAST for 30+ years

### 3. Wildcard Strategies Demo (Best for: Learning & comparison)

```bash
cd practice/simulate/dna_wildcard_strategies_demo
go run wildcard_comparison_demo.go
```

**Use when:**
- Evaluating different approaches
- Learning about algorithms
- Need to choose best strategy for your use case

**Shows:**
- 5 different wildcard strategies
- Pros/cons of each approach
- When to use which strategy

### 4. Original String Search (Best for: Single sequence, educational)

```bash
cd practice/simulate/dna_string_search
go run dna_string_search.go
```

**Features:**
- Naive O(n×m)
- KMP O(n+m)
- Sliding Window O(n×m)

## 📊 Algorithm Comparison

| Algorithm | Time Complexity | Space | Wildcards | Best For |
|-----------|----------------|-------|-----------|----------|
| **Naive** | O(N×M×Q) | O(1) | Unlimited | Small sequences |
| **KMP** | O(N+M) | O(M) | Unlimited | Single exact match |
| **N-gram** | O(2^W + H×Q) | O(N×M×k) | Limited | Repeated queries |
| **Seed-Extend** | O(H×Q) | O(N×M×k) | Unlimited | Large databases |

Where:
- N = number of sequences
- M = average sequence length
- Q = query length
- W = number of wildcards
- H = number of candidate matches
- k = k-mer size

## 🎯 Which Should I Use?

### Use **N-gram Search** if:
✅ You have a fixed database you'll query many times
✅ Wildcards are limited (1-3 per query)
✅ You want O(1) seed lookups
✅ Memory for index is acceptable

### Use **Seed-and-Extend** if:
✅ Your database is very large (genomic scale)
✅ Queries have many scattered wildcards
✅ You want battle-tested BLAST algorithm
✅ Index build time is one-time cost

### Use **Original String Search** if:
✅ Searching single sequences
✅ Learning about algorithms
✅ No database/index needed
✅ Educational purposes

## 🧬 Wildcard Support (IUPAC Codes)

All implementations support these wildcard nucleotides:

| Code | Matches | Meaning |
|------|---------|---------|
| A, C, G, T | Exact base | Standard nucleotides |
| R | A or G | puRine |
| Y | C or T | pYrimidine |
| M | A or C | aMino |
| K | G or T | Keto |
| W | A or T | Weak |
| S | C or G | Strong |
| B | C, G, T | not A |
| D | A, G, T | not C |
| H | A, C, T | not G |
| V | A, C, G | not T |
| N | A, C, G, T | aNy base |

## 📈 Performance Examples

**Query:** `ATTR` (R = A or G)
**Database:** 1,000,000 bases across 100 sequences

| Method | Index Build | Query Time | Total |
|--------|-------------|------------|-------|
| Naive | 0ms | 37ms | **37ms** |
| N-gram | 47ms | 2.7ms | **49.7ms** |
| Seed-Extend | 47ms | 2.7ms | **49.7ms** |

**For 100 queries:**
- Naive: 37ms × 100 = **3,700ms**
- N-gram: 47ms + 2.7ms × 100 = **317ms** (11.7× faster)
- Seed-Extend: 47ms + 2.7ms × 100 = **317ms** (11.7× faster)

## 📚 Documentation

Each folder contains detailed documentation:

- **README.md** - Implementation details and usage
- **BENEFITS.md** - Performance analysis
- **SEED_EXTEND.md** - Algorithm explanation
- **WILDCARD_STRATEGIES.md** - Strategy comparison

## 🔧 Requirements

All programs require:
- Go 1.16+ (for Go implementations)
- GCC/Clang (for C++ implementations)

No external dependencies needed!

## 🎓 Learning Path

**Beginner:**
1. Start with [dna_string_search/](../dna_string_search/) - Learn basic algorithms
2. Read about KMP and understand LPS table

**Intermediate:**
3. Try [dna_ngram_search_demo/](../dna_ngram_search_demo/) - Learn indexing
4. Understand wildcard expansion technique

**Advanced:**
5. Study [dna_seed_extend_demo/](../dna_seed_extend_demo/) - Production algorithm
6. Read [WILDCARD_STRATEGIES.md](../dna_wildcard_strategies_demo/WILDCARD_STRATEGIES.md) - Compare all approaches

## 🤝 Related Topics

- **String Matching:** Boyer-Moore, Rabin-Karp, Aho-Corasick
- **Bioinformatics:** BLAST, Smith-Waterman, Needleman-Wunsch
- **Indexing:** Suffix Arrays, Suffix Trees, FM-Index
- **Approximate Matching:** Edit Distance, Hamming Distance

## 📝 Summary

| Folder | Purpose | Run Command |
|--------|---------|-------------|
| `dna_string_search/` | Original algorithms | `go run dna_string_search.go` |
| `dna_ngram_search_demo/` | N-gram indexing | `go run dna_ngram_search.go` |
| `dna_seed_extend_demo/` | BLAST-style search | `go run dna_seed_extend.go` |
| `dna_wildcard_strategies_demo/` | Strategy comparison | `go run wildcard_comparison_demo.go` |

**All folders are independent and can be run with `package main`!** 🎉
