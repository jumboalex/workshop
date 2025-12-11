package main

import (
	"fmt"
	"strings"
)

/*
// Standard bases
Map.entry('A', List.of('A')),
Map.entry('C', List.of('C')),
Map.entry('G', List.of('G')),
Map.entry('T', List.of('T')),
// Degenerate bases
Map.entry('R', List.of('A', 'G')),
Map.entry('Y', List.of('C', 'T')),
Map.entry('M', List.of('A', 'C')),
Map.entry('K', List.of('G', 'T')),
Map.entry('W', List.of('A', 'T')),
Map.entry('S', List.of('C', 'G')),
Map.entry('B', List.of('C', 'G', 'T')),
Map.entry('D', List.of('A', 'G', 'T')),
Map.entry('H', List.of('A', 'C', 'T')),
Map.entry('V', List.of('A', 'C', 'G')),
Map.entry('N', List.of('A', 'C', 'G', 'T'))
);
*/
var patternCharMap map[byte]map[byte]struct{}

// demonstrateKMP shows a detailed step-by-step walkthrough of the KMP algorithm
func demonstrateKMP(sequence string, pattern string) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("DETAILED KMP ALGORITHM WALKTHROUGH")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Sequence: %s (length: %d)\n", sequence, len(sequence))
	fmt.Printf("Pattern:  %s (length: %d)\n", pattern, len(pattern))

	// Phase 1: Build LPS table
	fmt.Println("\n--- PHASE 1: Building LPS (Failure) Table ---")
	fmt.Println("LPS[i] = length of longest proper prefix which is also suffix")
	fmt.Println("This tells us where to restart when we have a mismatch\n")

	patLen := len(pattern)
	lps := make([]int, patLen)
	length := 0
	i := 1
	lps[0] = 0

	fmt.Printf("Step 0: lps[0] = 0 (always)\n")
	fmt.Printf("        Pattern: %s\n", pattern)
	fmt.Printf("        LPS:     [%d", lps[0])
	for k := 1; k < patLen; k++ {
		fmt.Printf(",_")
	}
	fmt.Println("]")

	step := 1
	for i < patLen {
		fmt.Printf("\nStep %d: i=%d, length=%d\n", step, i, length)
		fmt.Printf("        Comparing pattern[%d]='%c' with pattern[%d]='%c'\n", i, pattern[i], length, pattern[length])

		if patternCharsMatch(pattern[i], pattern[length]) {
			length++
			lps[i] = length
			fmt.Printf("        ✓ Match! length++ = %d, lps[%d] = %d\n", length, i, lps[i])
			i++
		} else {
			if length != 0 {
				fmt.Printf("        ✗ No match. length=%d, so jump to lps[%d] = %d\n", length, length-1, lps[length-1])
				length = lps[length-1]
			} else {
				lps[i] = 0
				fmt.Printf("        ✗ No match. length=0, so lps[%d] = 0, i++\n", i)
				i++
			}
		}

		fmt.Printf("        Pattern: %s\n", pattern)
		fmt.Printf("        LPS:     [")
		for k := 0; k < patLen; k++ {
			if k < i {
				if k > 0 {
					fmt.Printf(",")
				}
				fmt.Printf("%d", lps[k])
			} else {
				if k > 0 {
					fmt.Printf(",")
				}
				fmt.Printf("_")
			}
		}
		fmt.Println("]")
		step++
	}

	fmt.Println("\nFinal LPS table:", lps)
	fmt.Println("Interpretation:")
	for k := 0; k < patLen; k++ {
		if lps[k] > 0 {
			fmt.Printf("  lps[%d]=%d: If mismatch after position %d, restart at pattern[%d]\n", k, lps[k], k, lps[k])
		} else {
			fmt.Printf("  lps[%d]=%d: If mismatch after position %d, restart from beginning\n", k, lps[k], k)
		}
	}

	// Phase 2: Search using KMP
	fmt.Println("\n--- PHASE 2: Searching with KMP ---")
	fmt.Println("Key: i never goes backward, only j resets using LPS table\n")

	seqLen := len(sequence)
	i = 0
	j := 0
	step = 1

	for i < seqLen {
		fmt.Printf("Step %d: i=%d, j=%d\n", step, i, j)

		// Visual representation
		fmt.Printf("        Sequence: %s\n", sequence)
		fmt.Printf("                  %s^ (i=%d)\n", strings.Repeat(" ", i), i)
		fmt.Printf("        Pattern:  %s%s\n", strings.Repeat(" ", i-j), pattern)
		fmt.Printf("                  %s^ (j=%d)\n", strings.Repeat(" ", i), j)

		if j >= patLen {
			fmt.Printf("        ✓✓✓ MATCH FOUND at position %d!\n", i-patLen)
			break
		}

		seqChar := sequence[i]
		patChar := pattern[j]

		fmt.Printf("        Comparing sequence[%d]='%c' with pattern[%d]='%c'\n", i, seqChar, j, patChar)

		if _, ok := patternCharMap[patChar][seqChar]; ok {
			fmt.Printf("        ✓ Match! i++ and j++\n")
			i++
			j++
		} else {
			if j != 0 {
				fmt.Printf("        ✗ Mismatch! j=%d, so use lps[%d]=%d (NO backtracking of i!)\n", j, j-1, lps[j-1])
				j = lps[j-1]
			} else {
				fmt.Printf("        ✗ Mismatch! j=0, so just i++\n")
				i++
			}
		}
		fmt.Println()
		step++

		// Safety limit for demonstration
		if step > 50 {
			fmt.Println("        ... (truncated for brevity)")
			break
		}
	}

	if j == patLen {
		fmt.Printf("\n✓✓✓ MATCH FOUND at position %d!\n", i-patLen)
	} else if step <= 50 {
		fmt.Println("\n✗ No match found")
	}

	fmt.Println(strings.Repeat("=", 70))
}

// explainLPS provides a detailed explanation of what the LPS table means
func explainLPS(pattern string) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("UNDERSTANDING THE LPS (Longest Proper Prefix-Suffix) TABLE")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Pattern: %s\n\n", pattern)

	// Build the LPS table
	patLen := len(pattern)
	lps := buildKMPTable(pattern)

	fmt.Println("What is LPS?")
	fmt.Println("LPS[i] = Length of the longest PROPER prefix that is also a suffix")
	fmt.Println("         for the substring pattern[0...i]")
	fmt.Println()
	fmt.Println("A PROPER prefix/suffix means we exclude the string itself.")
	fmt.Println("Example: For 'ABAB':")
	fmt.Println("  - Proper prefixes: '', 'A', 'AB', 'ABA'")
	fmt.Println("  - Proper suffixes: '', 'B', 'AB', 'BAB'")
	fmt.Println("  - Common: '', 'AB' → longest is 'AB' with length 2")
	fmt.Println()

	// Show the LPS table
	fmt.Println("LPS Table for pattern:", pattern)
	fmt.Println()
	fmt.Print("Index:   ")
	for i := 0; i < patLen; i++ {
		fmt.Printf("%3d ", i)
	}
	fmt.Println()

	fmt.Print("Pattern: ")
	for i := 0; i < patLen; i++ {
		fmt.Printf("%3c ", pattern[i])
	}
	fmt.Println()

	fmt.Print("LPS:     ")
	for i := 0; i < patLen; i++ {
		fmt.Printf("%3d ", lps[i])
	}
	fmt.Println()
	fmt.Println()

	// Explain each position in detail
	fmt.Println("Detailed Explanation for Each Position:")
	fmt.Println(strings.Repeat("-", 70))

	for i := 0; i < patLen; i++ {
		substring := pattern[0 : i+1]
		fmt.Printf("\nPosition %d: pattern[0...%d] = \"%s\"\n", i, i, substring)

		if lps[i] == 0 {
			fmt.Printf("  LPS[%d] = 0\n", i)
			fmt.Println("  → No proper prefix matches a suffix")
		} else {
			prefix := substring[0:lps[i]]
			suffix := substring[len(substring)-lps[i]:]
			fmt.Printf("  LPS[%d] = %d\n", i, lps[i])
			fmt.Printf("  → Longest matching prefix/suffix: \"%s\"\n", prefix)
			fmt.Printf("  → Prefix: \"%s\" (first %d chars)\n", prefix, lps[i])
			fmt.Printf("  → Suffix: \"%s\" (last %d chars)\n", suffix, lps[i])
			fmt.Printf("  → They match!\n")
		}

		// Show visual representation
		fmt.Printf("\n  Visual:\n")
		fmt.Printf("  Pattern: %s\n", substring)
		if lps[i] > 0 {
			// Show prefix
			fmt.Printf("  Prefix:  ")
			for j := 0; j < len(substring); j++ {
				if j < lps[i] {
					fmt.Printf("%c", substring[j])
				} else {
					fmt.Printf("-")
				}
			}
			fmt.Println()

			// Show suffix
			fmt.Printf("  Suffix:  ")
			for j := 0; j < len(substring); j++ {
				if j < len(substring)-lps[i] {
					fmt.Printf("-")
				} else {
					fmt.Printf("%c", substring[j])
				}
			}
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Why is this useful?")
	fmt.Println()
	fmt.Println("When we have a MISMATCH at position j during search:")
	fmt.Println("  1. We've already matched pattern[0...j-1] with the sequence")
	fmt.Println("  2. LPS[j-1] tells us the longest prefix that matches the suffix")
	fmt.Println("  3. We can skip ahead: j = LPS[j-1]")
	fmt.Println("  4. This avoids re-checking characters we already know match!")
	fmt.Println()
	fmt.Println("Example: If we matched 'GATT' but fail at position 4:")
	fmt.Printf("  - We've matched: \"%s\"\n", pattern[0:min(4, patLen)])
	if 4 <= patLen && lps[min(3, patLen-1)] > 0 {
		fmt.Printf("  - LPS[3] = %d means the first %d chars match the last %d chars\n",
			lps[3], lps[3], lps[3])
		fmt.Printf("  - So we can restart at position %d instead of 0!\n", lps[3])
	}
	fmt.Println(strings.Repeat("=", 70))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	patternCharMap = make(map[byte]map[byte]struct{})
	patternCharMap['A'] = map[byte]struct{}{'A': {}}
	patternCharMap['C'] = map[byte]struct{}{'C': {}}
	patternCharMap['G'] = map[byte]struct{}{'G': {}}
	patternCharMap['T'] = map[byte]struct{}{'T': {}}
	patternCharMap['R'] = map[byte]struct{}{'A': {}, 'G': {}}
	patternCharMap['Y'] = map[byte]struct{}{'C': {}, 'T': {}}
	patternCharMap['M'] = map[byte]struct{}{'A': {}, 'C': {}}
	patternCharMap['K'] = map[byte]struct{}{'G': {}, 'T': {}}
	patternCharMap['W'] = map[byte]struct{}{'A': {}, 'T': {}}
	patternCharMap['S'] = map[byte]struct{}{'C': {}, 'G': {}}
	patternCharMap['B'] = map[byte]struct{}{'C': {}, 'G': {}, 'T': {}}
	patternCharMap['D'] = map[byte]struct{}{'A': {}, 'G': {}, 'T': {}}
	patternCharMap['H'] = map[byte]struct{}{'A': {}, 'C': {}, 'T': {}}
	patternCharMap['V'] = map[byte]struct{}{'A': {}, 'C': {}, 'G': {}}
	patternCharMap['N'] = map[byte]struct{}{'A': {}, 'C': {}, 'G': {}, 'T': {}}

	sequences := []string{
		"GATTACA", "GATTG",
	}
	pattern := "GATTR"

	// First, explain the LPS table in detail
	explainLPS("GATTR")
	explainLPS("AAAA")

	// Then show detailed KMP walkthroughs
	demonstrateKMP("GATGACA", "GATTR")

	fmt.Println("\n\n=== Testing All Three Algorithms ===")
	fmt.Printf("Pattern: %s\n\n", pattern)

	// Test with all three algorithms
	for _, seq := range sequences {
		naiveResult := searchDNASequenceNaive(seq, pattern)
		kmpResult := searchDNASequenceKMP(seq, pattern)
		slidingWindowResult := searchDNASequenceSlidingWindow(seq, pattern)
		fmt.Printf("Sequence: %s\n", seq)
		fmt.Printf("  Naive algorithm:          %v\n", naiveResult)
		fmt.Printf("  KMP algorithm:            %v\n", kmpResult)
		fmt.Printf("  Sliding Window algorithm: %v\n", slidingWindowResult)
		fmt.Println()
	}

	// Additional test cases to demonstrate both algorithms
	fmt.Println("=== Additional Test Cases ===")
	testCases := []struct {
		seq string
		pat string
		desc string
	}{
		{"AAAAAAAT", "AAAA", "Pattern with repeats (where KMP shines)"},
		{"AAAAG", "AAAR", "Pattern with wildcards"},
		{"GATTAGA", "GATTNR", "Complex pattern"},
		{"AAATTTGGG", "CCCC", "No match"},
	}

	ngramSize := 4
	for _, tc := range testCases {
		naiveResult := searchDNASequenceNaive(tc.seq, tc.pat)
		kmpResult := searchDNASequenceKMP(tc.seq, tc.pat)
		slidingWindowResult := searchDNASequenceSlidingWindow(tc.seq, tc.pat)
		ngramResult := searchDNASequenceNGram(tc.seq, tc.pat, ngramSize)
		fmt.Printf("%s\n", tc.desc)
		fmt.Printf("  Seq: %s, Pat: %s\n", tc.seq, tc.pat)
		fmt.Printf("  Naive: %v, KMP: %v, Sliding Window: %v, N-gram (n=%d): %v\n\n",
			naiveResult, kmpResult, slidingWindowResult, ngramSize, ngramResult)
	}

	// Demonstrate N-gram search with detailed example
	fmt.Println("\n=== N-Gram Search Detailed Example ===")
	fmt.Printf("N-gram size: %d\n\n", ngramSize)

	exampleSeq := "GATTACAGATTG"
	examplePat := "GATTR"

	fmt.Printf("Sequence: %s\n", exampleSeq)
	fmt.Printf("Pattern:  %s\n\n", examplePat)

	// Show n-grams from sequence
	fmt.Printf("4-grams from sequence:\n")
	for i := 0; i <= len(exampleSeq)-ngramSize; i++ {
		ngram := exampleSeq[i : i+ngramSize]
		fmt.Printf("  Position %d: %s\n", i, ngram)
	}

	fmt.Printf("\nPattern n-gram (first %d chars): %s\n", ngramSize, examplePat[0:ngramSize])
	fmt.Printf("Search result: %v\n\n", searchDNASequenceNGram(exampleSeq, examplePat, ngramSize))

	// Test: Can we find GATTA in ATTAGATT?
	fmt.Println("=== Special Test: Understanding Boundaries ===")
	testSeq1 := "ATTAGATT"
	testPat1 := "GATTA"

	fmt.Printf("Test 1 - Sequence: %s (length %d)\n", testSeq1, len(testSeq1))
	fmt.Printf("         Pattern:  %s (length %d)\n", testPat1, len(testPat1))
	fmt.Println("\nChecking all positions:")
	for i := 0; i <= len(testSeq1)-len(testPat1); i++ {
		window := testSeq1[i : i+len(testPat1)]
		fmt.Printf("  Position %d: %s vs %s - ", i, window, testPat1)
		if window == testPat1 {
			fmt.Println("✓ MATCH")
		} else {
			fmt.Println("✗ no match")
		}
	}

	fmt.Println("\nAlgorithm results:")
	fmt.Printf("  Naive:          %v\n", searchDNASequenceNaive(testSeq1, testPat1))
	fmt.Printf("  KMP:            %v\n", searchDNASequenceKMP(testSeq1, testPat1))
	fmt.Printf("  Sliding Window: %v\n", searchDNASequenceSlidingWindow(testSeq1, testPat1))
	fmt.Printf("  N-gram (n=4):   %v\n", searchDNASequenceNGram(testSeq1, testPat1, 4))

	fmt.Println("\n--- Now with sequence that DOES contain GATTA ---")
	testSeq2 := "ATTAGATTA"  // Added one more 'A' at the end
	fmt.Printf("Test 2 - Sequence: %s (length %d)\n", testSeq2, len(testSeq2))
	fmt.Printf("         Pattern:  %s (length %d)\n", testPat1, len(testPat1))
	fmt.Println("\nVisual:")
	fmt.Println("  A T T A G A T T A")
	fmt.Println("            G A T T A")
	fmt.Println("          (match at position 4)")

	fmt.Println("\nAlgorithm results:")
	fmt.Printf("  Naive:          %v\n", searchDNASequenceNaive(testSeq2, testPat1))
	fmt.Printf("  KMP:            %v\n", searchDNASequenceKMP(testSeq2, testPat1))
	fmt.Printf("  Sliding Window: %v\n", searchDNASequenceSlidingWindow(testSeq2, testPat1))
	fmt.Printf("  N-gram (n=4):   %v\n", searchDNASequenceNGram(testSeq2, testPat1, 4))
}

// buildKMPTable builds the failure function (partial match table) for KMP algorithm
// This handles wildcard patterns using IUPAC codes
func buildKMPTable(pattern string) []int {
	patLen := len(pattern)
	lps := make([]int, patLen) // longest proper prefix which is also suffix
	length := 0                 // length of the previous longest prefix suffix
	i := 1

	lps[0] = 0 // lps[0] is always 0

	for i < patLen {
		// Check if current pattern chars can match (considering wildcards)
		if patternCharsMatch(pattern[i], pattern[length]) {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
	return lps
}

// patternCharsMatch checks if two pattern characters can potentially match
// This is used for building the KMP table
func patternCharsMatch(p1, p2 byte) bool {
	// Get the possible bases for each pattern character
	bases1, ok1 := patternCharMap[p1]
	bases2, ok2 := patternCharMap[p2]

	if !ok1 || !ok2 {
		return false
	}

	// Check if there's any overlap in possible bases
	for base := range bases1 {
		if _, exists := bases2[base]; exists {
			return true
		}
	}
	return false
}

// searchDNASequenceNaive uses a naive string matching algorithm
// Time Complexity: O(n*m) where n is sequence length, m is pattern length
func searchDNASequenceNaive(sequence string, pattern string) bool {
	seqLen := len(sequence)
	patLen := len(pattern)

	i := 0
	j := 0

	for i < seqLen && j < patLen {
		seqChar := sequence[i]
		patChar := pattern[j]

		if _, ok := patternCharMap[patChar][seqChar]; ok {
			i++
			j++
		} else {
			i = i - j + 1 // Backtrack - this is inefficient
			j = 0
		}
	}
	return j == patLen
}

// searchDNASequenceKMP uses the KMP (Knuth-Morris-Pratt) algorithm
// Time Complexity: O(n+m) where n is sequence length, m is pattern length
func searchDNASequenceKMP(sequence string, pattern string) bool {
	seqLen := len(sequence)
	patLen := len(pattern)

	if patLen == 0 {
		return true
	}
	if seqLen == 0 {
		return false
	}

	// Build KMP failure table
	lps := buildKMPTable(pattern)

	i := 0 // index for sequence
	j := 0 // index for pattern

	for i < seqLen {
		seqChar := sequence[i]
		patChar := pattern[j]

		// Check if pattern character matches sequence character
		if _, ok := patternCharMap[patChar][seqChar]; ok {
			i++
			j++
		} else {
			// Mismatch after j matches
			if j != 0 {
				// Use KMP table to avoid redundant comparisons
				j = lps[j-1]
			} else {
				// No match at all, move to next sequence character
				i++
			}
		}

		// Found a complete match
		if j == patLen {
			return true
		}
	}

	return false
}

// searchDNASequenceSlidingWindow uses a sliding window approach
// Time Complexity: O(n*m) but simpler and often faster in practice due to early exit
func searchDNASequenceSlidingWindow(sequence string, pattern string) bool {
	seqLen := len(sequence)
	patLen := len(pattern)

	// Edge cases
	if patLen == 0 {
		return true
	}
	if seqLen < patLen {
		return false
	}

	// Slide the window of size patLen across the sequence
	for i := 0; i <= seqLen-patLen; i++ {
		matched := true

		// Check if all characters in the current window match the pattern
		for j := range patLen {
			seqChar := sequence[i+j]
			patChar := pattern[j]

			// Check if pattern character matches sequence character (considering wildcards)
			if _, ok := patternCharMap[patChar][seqChar]; !ok {
				matched = false
				break // Early exit on first mismatch
			}
		}

		if matched {
			return true
		}
	}

	return false
}

// searchDNASequenceNGram uses n-gram indexing for fast searching
// Time Complexity: O(n + m) for preprocessing, O(k*m) for search where k is # of n-gram matches
// Space Complexity: O(n) for the n-gram index
func searchDNASequenceNGram(sequence string, pattern string, n int) bool {
	seqLen := len(sequence)
	patLen := len(pattern)

	// Edge cases
	if patLen == 0 {
		return true
	}
	if seqLen < patLen {
		return false
	}
	if n <= 0 || n > patLen {
		// If n is invalid, fall back to sliding window
		return searchDNASequenceSlidingWindow(sequence, pattern)
	}

	// Step 1: Build n-gram index from sequence
	// Map from n-gram -> list of positions where it appears
	ngramIndex := make(map[string][]int)

	for i := 0; i <= seqLen-n; i++ {
		ngram := sequence[i : i+n]
		ngramIndex[ngram] = append(ngramIndex[ngram], i)
	}

	// Step 2: Extract first n-gram from pattern
	// Note: For patterns with wildcards, we need to consider all possible n-grams
	patternNGram := pattern[0:n]

	// Step 3: Check if this n-gram has wildcards
	// If it does, we need to generate all possible concrete n-grams
	hasWildcard := false
	for i := 0; i < n; i++ {
		if pattern[i] != 'A' && pattern[i] != 'C' && pattern[i] != 'G' && pattern[i] != 'T' {
			hasWildcard = true
			break
		}
	}

	// Step 4: Find candidate positions using n-gram index
	var candidatePositions []int

	if hasWildcard {
		// For wildcards, check all n-grams in the index
		for ngram, positions := range ngramIndex {
			if matchesNGramPattern(ngram, patternNGram) {
				candidatePositions = append(candidatePositions, positions...)
			}
		}
	} else {
		// Direct lookup for exact n-gram
		if positions, ok := ngramIndex[patternNGram]; ok {
			candidatePositions = positions
		}
	}

	// Step 5: Verify full pattern match at each candidate position
	for _, pos := range candidatePositions {
		if pos+patLen > seqLen {
			continue
		}

		matched := true
		for j := 0; j < patLen; j++ {
			seqChar := sequence[pos+j]
			patChar := pattern[j]

			if _, ok := patternCharMap[patChar][seqChar]; !ok {
				matched = false
				break
			}
		}

		if matched {
			return true
		}
	}

	return false
}

// matchesNGramPattern checks if an n-gram matches a pattern with wildcards
func matchesNGramPattern(ngram string, pattern string) bool {
	if len(ngram) != len(pattern) {
		return false
	}

	for i := 0; i < len(ngram); i++ {
		if _, ok := patternCharMap[pattern[i]][ngram[i]]; !ok {
			return false
		}
	}
	return true
}

// searchDNASequence is the default function (uses KMP for efficiency)
func searchDNASequence(sequence string, pattern string) bool {
	return searchDNASequenceKMP(sequence, pattern)
}

func searchAgain(sequence string, pattern string) bool {
	seqLen := len(sequence)
	patLen := len(pattern)

	if patLen == 0 {
		return true
	}
	if seqLen < patLen {
		return false
	}

	for i := 0; i <= seqLen-patLen; i++ {
		matched := true

		for j := 0; j < patLen; j++ {
			seqChar := sequence[i+j]
			patChar := pattern[j]
			if _, ok := patternCharMap[patChar][seqChar]; !ok {
				matched = false
				break
			}
		}

		if matched {
			return true
		}

	}
	return false
}