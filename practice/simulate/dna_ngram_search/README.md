# DNA N-Gram Search - Go vs C++ Implementation

## Overview

This implements a DNA sequence search engine using n-gram indexing with **full support** for IUPAC degenerate bases (wildcards).

## ✅ Current Status

The **Go implementation** (`dna_ngram_search.go`) includes **wildcard expansion** and correctly handles all test cases!

The **C++ implementation** (`dna_ngram_search.cpp`) has the original logic without wildcard expansion (for comparison).

## Algorithm

### 1. Preprocessing (Index Building)
- Extract all n-grams from each sequence
- Build inverted index: `n-gram -> list of sequences`

Example with n=4:
```
GATTACA → [GATT, ATTA, TTAC, TACA]
GATTG   → [GATT, ATTG]

Index:
GATT → [GATTACA, GATTG]
ATTA → [GATTACA]
TTAC → [GATTACA]
TACA → [GATTACA]
ATTG → [GATTG]
```

### 2. Query Search
- Extract n-grams from query
- Lookup each n-gram in index
- Intersect results to get candidates
- Verify each candidate with full pattern matching

## Test Results (Go Implementation with Wildcard Expansion)

### ✅ All Tests Pass!

**Query: "GATT"**
- N-grams: [GATT]
- Expanded: [GATT] (no wildcards)
- Result: **[GATTACA, GATTG]** ✓

**Query: "ATTACA"**
- N-grams: [ATTA, TTAC, TACA]
- Expanded: All exact (no wildcards)
- Result: **[GATTACA]** ✓

**Query: "GATTR"** (R = A or G) - **NOW WORKS!** ✓
- N-grams: [GATT, ATTR]
- Expanded: [GATT], [ATTA, ATTG]
- Lookup: Both ATTA and ATTG found in index
- Result: **[GATTACA, GATTG]** ✓✓✓

**Query: "GATTM"** (M = A or C) - **NOW WORKS!** ✓
- N-grams: [GATT, ATTM]
- Expanded: [GATT], [ATTA, ATTC]
- Result: **[GATTACA]** ✓✓✓

**Query: "GATTRR"** (RR = no matches)
- N-grams: [GATT, ATTR, TTRR]
- Expanded: Multiple variants, but TTRR variants don't exist
- Result: **[]** ✓ (correctly empty)

## Solutions for Wildcard Support

### Option 1: Expand Wildcards During Indexing
When building index, expand wildcards:
```
ATTR could match: ATTA, ATTG
Add both to index
```
❌ Not applicable here - sequences don't have wildcards

### Option 2: Expand Wildcards During Query
When query has wildcards, generate all possible concrete n-grams:
```
ATTR (R=A or G) → [ATTA, ATTG]
Look up both in index
```
✓ This would work!

### Option 3: Wildcard-Aware Matching in Lookup
Instead of exact string match, use wildcard matching:
```
For n-gram "ATTR" with wildcard:
  Check all index entries
  Match "ATTR" pattern against each
```
✓ This would work but slower

### Option 4: Fallback to Direct Search
If query has wildcards, skip n-gram optimization:
```
if (query has wildcards in n-grams):
    use direct search on all sequences
```
✓ Simple but loses n-gram benefits

## Current Implementation

### Go Implementation (`dna_ngram_search.go`)
- ✅ Build n-gram index correctly
- ✅ Handle exact pattern matching
- ✅ **Expand wildcards in n-grams** ← FIXED!
- ✅ Intersect candidates properly
- ✅ Verify with wildcard support in final step
- ✅ **All test cases pass!**

### C++ Implementation (`dna_ngram_search.cpp`)
- ✅ Build n-gram index correctly
- ✅ Handle exact pattern matching
- ✅ Intersect candidates properly
- ✅ Verify with wildcard support in final step
- ⚠️ Don't handle wildcards in n-gram lookup phase (original version for comparison)

## Comparison: Go vs C++

### Similarities
- Same algorithm and logic
- Same test results
- Both handle IUPAC wildcards in verification step
- Both fail on wildcard n-gram lookup

### Differences

| Aspect | Go | C++ |
|--------|-----|-----|
| **Syntax** | Cleaner, more concise | More verbose |
| **Maps** | `map[string]map[string]struct{}` | `unordered_map<string, unordered_set<string>>` |
| **Sets** | Use `map[T]struct{}` idiom | Native `unordered_set` |
| **Printing** | `fmt.Printf` | `cout <<` |
| **Iteration** | Range-based: `for k, v := range` | Iterator-based or range (C++11+) |
| **Performance** | Fast (compiled) | Slightly faster (more control) |
| **Compilation** | Very fast | Slower |

## Usage

### Go
```bash
cd /home/jumbo/workspace/workshop/practice/simulate/dna_ngram_search
go run dna_ngram_search.go
```

### C++
```bash
cd /home/jumbo/workspace/workshop/practice/simulate/dna_ngram_search
g++ -std=c++11 -o dna_ngram_cpp dna_ngram_search.cpp
./dna_ngram_cpp
```

## Files in This Directory

- **`dna_ngram_search.go`** - Go implementation with full wildcard support ✅
- **`dna_ngram_search.cpp`** - C++ implementation (original, for comparison)
- **`README.md`** - This file
- **`WILDCARD_FIX.md`** - Detailed explanation of how wildcards were fixed

## Key Achievement

✅ **Successfully implemented Option 2** - Wildcard expansion during query processing!

The Go implementation now:
1. Detects wildcards in n-grams
2. Expands each wildcard n-gram to all possible concrete variants
3. Looks up all variants in the index
4. Unions the results before intersection
5. Passes all test cases with wildcards!
