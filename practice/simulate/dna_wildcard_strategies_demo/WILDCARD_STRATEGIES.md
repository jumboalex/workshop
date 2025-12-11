# Efficient Wildcard Search Strategies

Beyond wildcard expansion, there are several other approaches to handle wildcards efficiently in DNA search.

## 🔍 Strategy Comparison

### **Strategy 1: Wildcard Expansion (Current Implementation)**

**How it works:**
```
Query: "ATTR" (R = A or G)
Expand: ["ATTA", "ATTG"]
Lookup both in index
```

**Pros:**
- ✅ Exact matches in index (fast lookups)
- ✅ Reuses existing index structure
- ✅ Simple to implement

**Cons:**
- ❌ Exponential explosion with multiple wildcards
- ❌ "NNNN" = 4^4 = 256 variants!
- ❌ Redundant lookups

**Complexity:**
- Time: O(4^w × log n) where w = wildcards in n-gram
- Space: O(1) extra (just query processing)

---

### **Strategy 2: Bit Vector Index**

**How it works:**
Store each base as 2 bits (A=00, C=01, G=10, T=11)
```
Sequence: GATT
Binary:   10 00 11 11

Wildcard R (A or G) = 00 or 10
Create bit mask for matching
```

**Implementation:**
```go
type BitVectorIndex struct {
    // Each sequence encoded as bit array
    sequences map[string]uint64
    // Each n-gram maps to bitmask
    ngramBits map[string]uint64
}

// Query with wildcards
func (b *BitVectorIndex) Search(query string) {
    queryBits := encode(query)
    wildcardMask := createWildcardMask(query)

    // Bitwise AND for fast matching
    for seq, bits := range b.sequences {
        if (bits & wildcardMask) == queryBits {
            // Match!
        }
    }
}
```

**Pros:**
- ✅ Very fast: single CPU instruction (bitwise AND)
- ✅ No explosion with multiple wildcards
- ✅ Cache-friendly (compact representation)

**Cons:**
- ❌ Only works for exact-length patterns
- ❌ Complex to implement
- ❌ Limited to sequences that fit in register (32-64 bases)

**Complexity:**
- Time: O(n) but with very small constant
- Space: O(n × m/4) (compressed)

**Use Case:** Perfect for fixed-length k-mer matching in genomics

---

### **Strategy 3: Trie/Suffix Tree with Wildcard Matching**

**How it works:**
Build a trie/suffix tree, traverse it with wildcard logic
```
        Root
       /  |  \
      A   C   G   T
     / \  |  / \  |
    ...

Query "ATR" (R = A|G):
- Follow A → T → both A and G branches
```

**Implementation:**
```go
type TrieNode struct {
    children map[byte]*TrieNode
    sequences []string // sequences containing this path
}

func (t *Trie) SearchWildcard(pattern string, pos int, node *TrieNode) {
    if pos == len(pattern) {
        return node.sequences
    }

    char := pattern[pos]
    if isWildcard(char) {
        results := []string{}
        // Follow all possible branches
        for base := range wildcardBases[char] {
            results = append(results,
                t.SearchWildcard(pattern, pos+1, node.children[base])...)
        }
        return results
    } else {
        return t.SearchWildcard(pattern, pos+1, node.children[char])
    }
}
```

**Pros:**
- ✅ Handles wildcards naturally during traversal
- ✅ No exponential expansion
- ✅ Prefix-based sharing saves space
- ✅ Works for variable-length patterns

**Cons:**
- ❌ More complex data structure
- ❌ Higher memory usage than hash table
- ❌ Slower than hash lookups for exact matches
- ❌ Build time is slower

**Complexity:**
- Time: O(m × 4^w) where w = wildcards
- Space: O(n × m) worst case

**Use Case:** Best when you have many overlapping sequences

---

### **Strategy 4: Automata-Based (DFA/NFA)**

**How it works:**
Build a deterministic/non-deterministic finite automaton
```
For pattern "ATR" (R = A|G):

State 0 --A--> State 1 --T--> State 2 --A--> Accept
                                        --G--> Accept
```

**Implementation:**
```go
type Automaton struct {
    states []State
    transitions map[State]map[byte]State
}

func (a *Automaton) Match(sequence string) bool {
    state := a.start
    for _, char := range sequence {
        state = a.transitions[state][char]
        if state == InvalidState {
            return false
        }
    }
    return state == AcceptState
}
```

**Pros:**
- ✅ Very efficient matching (linear scan)
- ✅ Handles complex wildcard patterns
- ✅ Can be compiled/optimized
- ✅ Good for regular expression-like patterns

**Cons:**
- ❌ Complex to build
- ❌ DFA can have state explosion
- ❌ Still requires scanning sequences

**Complexity:**
- Time: O(m) per sequence, O(n×m) total
- Space: O(4^m) for DFA, O(m) for NFA

**Use Case:** When wildcards form complex patterns (regex-like)

---

### **Strategy 5: Inverted Index with Wildcard Positions**

**How it works:**
Store separate indices for each wildcard pattern
```
Index 1: Exact n-grams (no wildcards)
  GATT → [Seq1, Seq2]
  ATTA → [Seq1]

Index 2: Position-specific wildcards
  GA** → [Seq1, Seq3]  (wildcard at position 2-3)
  **TT → [Seq1, Seq2]  (wildcard at position 0-1)

Query "GATTR":
  Lookup in exact: none
  Lookup in wildcard: GA** matches
  Verify candidates
```

**Implementation:**
```go
type WildcardIndex struct {
    exact map[string][]string
    // Key = pattern with wildcards at specific positions
    wildcardPatterns map[string]map[string][]string
}
```

**Pros:**
- ✅ Pre-computes common wildcard patterns
- ✅ Fast lookup for known patterns
- ✅ No runtime expansion needed

**Cons:**
- ❌ Huge index size (many possible patterns)
- ❌ Must predict which patterns will be queried
- ❌ Not flexible to arbitrary wildcards

**Complexity:**
- Time: O(1) lookup if pattern is indexed
- Space: O(n × m × 4^w) - HUGE!

**Use Case:** When wildcard positions are predictable and limited

---

### **Strategy 6: Seed-and-Extend (BLAST-like)**

**How it works:**
Find exact match "seeds", then extend with wildcards
```
Query: "GATTRACA" (R = A|G)

Step 1: Find exact seed "GATT"
  Index: GATT → [Seq1, Seq5, Seq9]

Step 2: Extend from seed with wildcard matching
  Seq1: GATT[A]ACA → matches R ✓
  Seq5: GATT[G]ACA → matches R ✓
  Seq9: GATT[C]XXX → no match ✗
```

**Implementation:**
```go
func SeedAndExtend(query string, index NgramIndex) []string {
    // Find longest exact substring (seed)
    seed := findLongestExactSubstring(query)

    // Lookup seed in index
    candidates := index[seed]

    // Extend with full wildcard matching
    results := []string{}
    for _, seq := range candidates {
        if extendMatch(seq, query, seedPosition) {
            results = append(results, seq)
        }
    }
    return results
}
```

**Pros:**
- ✅ Combines benefits of index lookup and flexible matching
- ✅ Reduces search space significantly
- ✅ Works well with few wildcards
- ✅ Used by BLAST (proven effective)

**Cons:**
- ❌ If no good seed, degrades to linear search
- ❌ Requires heuristics to choose seed
- ❌ May miss matches if seed doesn't match

**Complexity:**
- Time: O(k × m) where k = candidates from seed
- Space: O(n × m) for index

**Use Case:** Real-world DNA alignment (BLAST, BLAT)

---

### **Strategy 7: Bloom Filter for Quick Rejection**

**How it works:**
Use space-efficient probabilistic data structure
```
Bloom Filter: Does sequence CONTAIN n-gram?
- Insert all n-grams from all sequences
- Query: Check if pattern n-grams exist
- If NO → definitely no match
- If YES → maybe match (verify)
```

**Implementation:**
```go
type BloomIndex struct {
    filter *BloomFilter
    sequences []string
}

func (b *BloomIndex) Search(query string) []string {
    ngrams := extractNGrams(query)

    candidates := []string{}
    for _, seq := range b.sequences {
        maybeMatch := true
        for _, ngram := range ngrams {
            if !b.filter.Contains(seq, ngram) {
                maybeMatch = false
                break
            }
        }
        if maybeMatch {
            candidates = append(candidates, seq)
        }
    }

    // Verify candidates
    return verify(candidates, query)
}
```

**Pros:**
- ✅ Very space-efficient
- ✅ Fast negative checks
- ✅ Good for filtering before expensive operations

**Cons:**
- ❌ False positives (must verify)
- ❌ Still need to store sequences
- ❌ Not much better than regular index for wildcards

**Complexity:**
- Time: O(n) with small constant
- Space: O(n) but much smaller than full index

**Use Case:** First-pass filter for very large databases

---

### **Strategy 8: Hybrid: Index + Linear Scan with SIMD**

**How it works:**
Use index to narrow down + optimized linear scan
```
Step 1: Index lookup (filter to 1% of sequences)
Step 2: SIMD-optimized wildcard matching on candidates

SIMD = Single Instruction Multiple Data
Match 16 bases in parallel with CPU vectorization
```

**Implementation:**
```go
import "golang.org/x/sys/cpu"

func SIMDMatch(sequence string, pattern string) bool {
    // Use CPU vector instructions (AVX2/SSE)
    // Process 16/32 characters at once
    // Hardware-level parallelism
}

func HybridSearch(query string, index NgramIndex) []string {
    // Use index for coarse filtering
    candidates := indexLookup(query)

    // SIMD for fine-grained matching
    results := []string{}
    for _, seq := range candidates {
        if SIMDMatch(seq, query) {
            results = append(results, seq)
        }
    }
    return results
}
```

**Pros:**
- ✅ Extremely fast verification (10x+ speedup)
- ✅ Leverages modern CPU features
- ✅ Best of both worlds

**Cons:**
- ❌ Platform-specific code
- ❌ Complex implementation
- ❌ Requires understanding of CPU architectures

**Complexity:**
- Time: O(k × m/16) with SIMD speedup
- Space: Same as regular index

**Use Case:** High-performance production systems

---

## 📊 Comparison Table

| Strategy | Speed | Memory | Wildcards | Complexity | Best For |
|----------|-------|--------|-----------|------------|----------|
| **Expansion** | ★★★★☆ | ★★★★★ | Limited | ★★★★★ | Few wildcards |
| **Bit Vector** | ★★★★★ | ★★★★★ | Good | ★★☆☆☆ | Fixed-length k-mers |
| **Trie** | ★★★☆☆ | ★★★☆☆ | Excellent | ★★★☆☆ | Prefix searches |
| **Automaton** | ★★★★☆ | ★★★★☆ | Excellent | ★★☆☆☆ | Complex patterns |
| **Inverted Wildcard** | ★★★★★ | ★☆☆☆☆ | Limited | ★★★★☆ | Predictable patterns |
| **Seed-Extend** | ★★★★☆ | ★★★★☆ | Good | ★★★☆☆ | Real-world DNA |
| **Bloom Filter** | ★★★☆☆ | ★★★★★ | Good | ★★★★☆ | Huge databases |
| **Hybrid SIMD** | ★★★★★ | ★★★★☆ | Excellent | ★☆☆☆☆ | Production systems |

---

## 🎯 Recommendations

### For Your Use Case (DNA Search with IUPAC Codes)

**Current (Expansion):** ✅ Good choice!
- Works well for 1-2 wildcards
- Simple to implement
- Fast enough for most queries

**Upgrade Path:**

1. **For better performance:** → **Seed-and-Extend**
   - Use longest exact substring as seed
   - More robust with multiple wildcards
   - Industry-proven (BLAST uses this)

2. **For maximum speed:** → **Bit Vector + SIMD**
   - Hardware-accelerated matching
   - 10-100x faster verification
   - Worth it for high-throughput systems

3. **For complex patterns:** → **Automaton**
   - If users need regex-like patterns
   - More flexible than simple wildcards

### Quick Decision Tree

```
How many wildcards per query?
├─ 0-2: Current expansion ✓
├─ 3-5: Seed-and-extend
└─ 6+: Automaton or Trie

How many sequences?
├─ < 1,000: Current is fine ✓
├─ 1,000-1M: Consider seed-extend
└─ 1M+: SIMD + index

Need regex patterns?
├─ No: Current is fine ✓
└─ Yes: Build automaton

Have time to optimize?
├─ No: Keep current ✓
└─ Yes: Implement SIMD
```

---

## 💡 Conclusion

**Your current wildcard expansion is actually a great choice** for:
- ✅ Few wildcards (1-3 per query)
- ✅ Simple IUPAC codes
- ✅ Reasonable database size
- ✅ Easy to maintain

**Consider upgrading if:**
- Many wildcards (4+)
- Massive database (millions of sequences)
- Need maximum performance
- Complex query patterns

The best approach depends on your specific constraints and requirements! 🎯
