# Seed-and-Extend Algorithm for DNA Wildcard Search

## Overview

**Seed-and-extend** is the gold standard algorithm used by BLAST (Basic Local Alignment Search Tool) and other professional bioinformatics tools. It efficiently handles wildcard searches by using exact-match "seeds" as anchor points, then extending to verify the full pattern.

## How It Works

### Two-Phase Approach

```
Phase 1: SEED         Phase 2: EXTEND
┌──────────┐         ┌──────────┐
│ Find     │         │ Verify   │
│ Exact    │────────▶│ Full     │
│ Matches  │         │ Pattern  │
└──────────┘         └──────────┘
```

### Example Walkthrough

**Query:** `ATTR` (where R = A or G)

**Sequence:** `GATTACATTAGC`

#### Phase 1: Extract Seeds

```
Query: A T T R
       └─┬─┘ │
       Seed  Wildcard
       "ATT"
```

Seeds are the longest exact-match segments (no wildcards).

#### Phase 2: Find Seed Positions

```
Sequence: G A T T A C A T T A G C
Index:    0 1 2 3 4 5 6 7 8 9 10 11
              ^─^─^       ^─^─^
              Seed        Seed
              @pos 1      @pos 6
```

#### Phase 3: Extend and Verify

**Position 1:**
```
Sequence: G [A T T A] C A T T A G C
Query:      [A T T R]
Match:       ✓ ✓ ✓ ✓  (A matches R)
Result: MATCH at position 1
```

**Position 6:**
```
Sequence: G A T T A C [A T T A] G C
Query:                [A T T R]
Match:                 ✓ ✓ ✓ ✓  (A matches R)
Result: MATCH at position 6
```

## Algorithm Details

### 1. Index Building

```go
// For each sequence, index all k-mers
for i := 0; i <= len(sequence) - minSeedLen; i++ {
    seed := sequence[i : i + minSeedLen]
    index[seed] = append(index[seed], position{seqID, i})
}
```

**Example:** Sequence `GATTAC` with minSeedLen=3

```
Index:
  "GAT" → [pos 0]
  "ATT" → [pos 1]
  "TTA" → [pos 2]
  "TAC" → [pos 3]
```

### 2. Seed Extraction

```go
// Split query into non-wildcard segments
Query: "ATWATTR"
       └┬┘└──┬─┘
        │    └─── Seed: "ATT" at offset 3
        └──────── Seed: "AT" at offset 0 (too short if minSeedLen=3)
```

### 3. Seed Lookup

```go
// Look up each seed in the index
for each seed in seeds {
    matches := index[seed.sequence]
    for each match in matches {
        candidate_position = match.position - seed.offset
        candidates.add(candidate_position)
    }
}
```

### 4. Extension & Verification

```go
// For each candidate position, check full pattern
for each candidate in candidates {
    if fullPatternMatches(sequence, query, candidate) {
        results.add(candidate)
    }
}
```

## Complexity Analysis

| Operation | Time Complexity | Space Complexity |
|-----------|-----------------|------------------|
| **Index Building** | O(N × M) | O(N × M) |
| **Seed Extraction** | O(Q) | O(Q) |
| **Seed Lookup** | O(S) | O(1) |
| **Extension** | O(H × Q) | O(1) |
| **Total Query** | **O(S + H × Q)** | **O(N × M)** |

Where:
- N = number of sequences
- M = average sequence length
- Q = query length
- S = number of unique seeds
- H = number of seed hits (candidates)

## Performance Characteristics

### Best Case
```
Query: "ATTTTTTTG" (one small wildcard at end)
Seed: "ATTTTTT" (very specific)
Hits: Few matches → Fast extension
Time: O(Q) - nearly instant
```

### Average Case
```
Query: "ATRWATT" (2-3 wildcards)
Seeds: "AT", "ATT" (moderate specificity)
Hits: Moderate matches
Time: O(H × Q) where H << N
```

### Worst Case
```
Query: "NNNN" (all wildcards)
Seeds: None!
Fallback: Check every position
Time: O(N × M × Q) - same as naive
```

## Comparison with Other Approaches

| Approach | Query Time | Memory | Wildcards | Complexity |
|----------|-----------|---------|-----------|-----------|
| **Naive** | O(N×M×Q) | O(N×M) | Unlimited | Simple |
| **N-gram Expansion** | O(2^W × Q) | O(N×M) | Limited | Simple |
| **Seed-Extend** | **O(H×Q)** | **O(N×M)** | **Unlimited** | **Medium** |
| **Automaton** | O(N×M) | O(2^W×Q) | Unlimited | Complex |

Where W = number of wildcards

## Advantages

✅ **Proven Technology** - Used by BLAST for 30+ years
✅ **Handles Many Wildcards** - No exponential explosion
✅ **Good Performance** - Only checks promising positions
✅ **Reasonable Memory** - Just stores k-mer index
✅ **Natural for DNA** - Bioinformatics standard

## Disadvantages

❌ **Index Size** - Needs O(N×M) space for k-mer index
❌ **No Seeds = Slow** - Falls back to naive for all-wildcard queries
❌ **Parameter Tuning** - Need to choose minSeedLen carefully

## When to Use

### ✅ Use Seed-and-Extend When:
- Searching genomic databases (DNA/RNA sequences)
- Queries have mix of exact bases and wildcards
- Need to handle many wildcards (3+)
- Memory for k-mer index is acceptable
- Want proven, battle-tested algorithm

### ❌ Don't Use When:
- All queries are exact matches (use hash table)
- Queries are all wildcards (use naive or automaton)
- Memory is extremely limited (use expansion + compression)
- Need regex-like features (use automaton)

## Parameter Selection

### Minimum Seed Length (minSeedLen)

```
Too Small (minSeedLen = 2):
  Pro: More seeds found
  Con: Too many false hits, slow extension

Optimal (minSeedLen = 3-5):
  Pro: Good balance of specificity and sensitivity
  Con: None

Too Large (minSeedLen = 10+):
  Pro: Very few false hits
  Con: Might miss matches if seeds span wildcards
```

**Recommendation:** Use minSeedLen = 3 for DNA (4 bases), minSeedLen = 4-5 for proteins (20 amino acids)

## Real-World Example: BLAST

BLAST uses seed-and-extend with optimizations:

1. **Word Hits** - Seeds of length 11 (DNA) or 3 (protein)
2. **Two-Hit Method** - Requires 2 nearby seeds for extension
3. **Gapped Extension** - Handles insertions/deletions
4. **Score Matrices** - Uses substitution matrices for fuzzy matching

## Code Example

```go
// Create engine with minimum seed length of 3
engine := NewDNASeedExtendEngine(3)

// Add sequences to database
engine.AddSequence("seq1", "GATTACATTAGC")
engine.AddSequence("seq2", "CCGATTAGGATT")

// Search with wildcards
results := engine.Search("ATTR")  // R = A or G

// Results:
// seq1: [1, 6]  - matches at positions 1 and 6
// seq2: [3, 8]  - matches at positions 3 and 8
```

## Optimization Tips

### 1. Filter by Seed Count
```go
// Require at least 2 seeds for extension
if len(seeds) < 2 {
    return nil  // Query too ambiguous
}
```

### 2. Use Multiple Seed Sizes
```go
// Try 3-mers first, fall back to 2-mers if needed
seeds := extractSeeds(query, 3)
if len(seeds) == 0 {
    seeds = extractSeeds(query, 2)
}
```

### 3. Rank Candidates
```go
// Prioritize positions with multiple seed hits
candidates := rankByNumberOfSeedsHit(allCandidates)
for _, candidate := range candidates {
    if verify(candidate) {
        return candidate  // Stop on first match
    }
}
```

### 4. Early Termination
```go
// Stop extending if mismatch found
for i := 0; i < len(query); i++ {
    if !matches(query[i], sequence[pos+i]) {
        return false  // No need to check remaining positions
    }
}
```

## Variants

### Ungapped Seed-Extend
- Original BLAST algorithm
- Seeds must be continuous exact matches
- Fast but less sensitive

### Gapped Seed-Extend
- Modern BLAST (BLAST+)
- Allows insertions/deletions during extension
- More sensitive for divergent sequences

### Spaced Seeds
- Seeds can have wildcard positions: `A*T*G` (where * is "don't care")
- More sensitive for finding distant homologs
- Used in tools like PatternHunter

### Multiple Seeds
- Require 2+ nearby seeds before extending
- Reduces false positives
- Used in BLAST's two-hit method

## Further Reading

- **BLAST Paper**: Altschul et al. (1990) "Basic local alignment search tool"
- **PatternHunter**: Ma et al. (2002) "PatternHunter: faster and more sensitive homology search"
- **Spaced Seeds**: Choi et al. (2004) "Good spaced seeds for homology search"

## Summary

Seed-and-extend is the **best general-purpose algorithm** for DNA wildcard search:

| Metric | Rating | Notes |
|--------|--------|-------|
| Speed | ⭐⭐⭐⭐ | Fast for typical queries |
| Memory | ⭐⭐⭐ | Moderate k-mer index |
| Wildcards | ⭐⭐⭐⭐⭐ | Handles any number |
| Simplicity | ⭐⭐⭐ | Medium complexity |
| Proven | ⭐⭐⭐⭐⭐ | 30+ years in production |

**Bottom Line:** If BLAST uses it to search the entire human genome in seconds, it's good enough for your DNA search! 🧬
