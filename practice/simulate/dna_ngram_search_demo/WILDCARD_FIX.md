# How Our Original Wildcard Algorithm Fixed the N-Gram Search

## The Problem

In the original n-gram implementation, wildcard queries failed:

```
Query: "GATTR" (R = A or G)
N-grams: [GATT, ATTR]
Lookup: ATTR → [] (doesn't exist in index!)
Result: [] ❌ (WRONG - should find both sequences)
```

## The Solution: Wildcard Expansion

We leveraged our **original wildcard matching algorithm** to expand wildcard n-grams into all possible concrete variants!

### Key Components from Original Algorithm

#### 1. **patternCharMap** - The Foundation
```go
patternCharMap['R'] = map[byte]struct{}{'A': {}, 'G': {}}
patternCharMap['M'] = map[byte]struct{}{'A': {}, 'C': {}}
// ... etc
```

This map tells us which bases a wildcard can match.

#### 2. **Recursive Expansion Logic**
```go
func expandWildcardNGram(ngram string) []string {
    // Check if has wildcard
    // If yes, recursively expand each position
    // Generate all combinations
}
```

### How It Works

**Example: "ATTR" with wildcard R**

```
Position 0: A → [A]
Position 1: T → [T]
Position 2: T → [T]
Position 3: R → [A, G]  ← Wildcard!

Combinations:
A + T + T + A = "ATTA"
A + T + T + G = "ATTG"

Result: ["ATTA", "ATTG"]
```

Now we can look up BOTH variants in the index!

### Complete Flow

**Query: "GATTR"**

1. **Extract N-grams**
   ```
   [GATT, ATTR]
   ```

2. **Expand Wildcards** ← NEW STEP!
   ```
   GATT → [GATT]       (no wildcard)
   ATTR → [ATTA, ATTG] (R expanded)
   ```

3. **Lookup in Index**
   ```
   [GATT]       → {GATTACA, GATTG}
   [ATTA, ATTG] → {GATTACA, GATTG} (union of both)
   ```

4. **Intersect**
   ```
   {GATTACA, GATTG} ∩ {GATTACA, GATTG} = {GATTACA, GATTG}
   ```

5. **Verify** (using original matching logic)
   ```
   GATTACA: GATT[A] vs GATTR → A matches R ✓
   GATTG:   GATT[G] vs GATTR → G matches R ✓
   ```

6. **Result**
   ```
   [GATTACA, GATTG] ✓✓✓
   ```

## Test Results - FIXED! ✅

### Before (Original)
```
GATT   → [GATTACA, GATTG]  ✓
ATTACA → [GATTACA]          ✓
GATTR  → []                 ❌
GATTM  → []                 ❌
GATTRR → []                 ✓
```

### After (With Wildcard Expansion)
```
GATT   → [GATTACA, GATTG]  ✓
ATTACA → [GATTACA]          ✓
GATTR  → [GATTACA, GATTG]  ✓✓✓ FIXED!
GATTM  → [GATTACA]          ✓✓✓ FIXED!
GATTRR → []                 ✓ (correctly no match)
```

## Detailed Example: "GATTM"

**Query:** "GATTM" (M = A or C)

**Step 1: Expand**
```
GATT → [GATT]
ATTM → [ATTA, ATTC]  (M expands to A and C)
```

**Step 2: Lookup**
```
GATT → Index lookup → {GATTACA, GATTG}
ATTA → Index lookup → {GATTACA}
ATTC → Index lookup → {} (doesn't exist)

Union: {GATTACA} ∪ {} = {GATTACA}
```

**Step 3: Intersect**
```
{GATTACA, GATTG} ∩ {GATTACA} = {GATTACA}
```

**Step 4: Verify**
```
GATTACA: GATT[A]CA
         GATT[M]    → A matches M ✓
```

**Result:** [GATTACA] ✓

## Complexity Analysis

### Without Wildcard Expansion
- Time: O(k) where k = number of n-grams
- Lookups: k direct hash lookups
- **Problem:** Misses wildcard matches

### With Wildcard Expansion
- Time: O(k * 4^w) where w = wildcards per n-gram
- Lookups: k * (4^w) hash lookups (worst case with N wildcard)
- **Benefit:** Finds all wildcard matches correctly!

**Worst case:** N wildcard (matches all 4 bases)
```
"NNNN" → 4^4 = 256 variants
```

**Typical case:** 1-2 wildcards
```
"ATTR" (R = 2 bases) → 2 variants
"GATM" (M = 2 bases) → 2 variants
```

## Key Insights

1. **Reusability:** The wildcard matching logic from our original algorithm was **directly applicable** to n-gram expansion

2. **Modularity:** By separating the wildcard expansion into its own function, we made the code clean and maintainable

3. **Efficiency:** We only expand wildcards during query time (not during indexing), keeping the index compact

4. **Correctness:** The final verification step (using original matching) catches any false positives

## Code Changes

**Added 2 functions:**
```go
expandWildcardNGram()        // Detects and expands wildcards
expandWildcardHelper()       // Recursive expansion logic
```

**Modified 1 function:**
```go
Search()  // Now calls expandWildcardNGram() before lookup
```

**Result:** Full wildcard support with minimal code changes!

## Conclusion

✅ Our original wildcard algorithm was **essential** for fixing the n-gram search

✅ The `patternCharMap` structure provided the foundation for expansion

✅ The recursive matching logic translated perfectly to n-gram generation

✅ All test cases now pass correctly!

This is a great example of how **modular, well-designed algorithms** can be reused and extended to solve related problems! 🎉
