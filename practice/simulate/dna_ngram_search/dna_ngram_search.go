package main

import (
	"fmt"
	"strings"
)

// IUPAC nucleotide codes mapping
var patternCharMap map[byte]map[byte]struct{}

func initializePatternMap() {
	patternCharMap = make(map[byte]map[byte]struct{})
	// Standard bases
	patternCharMap['A'] = map[byte]struct{}{'A': {}}
	patternCharMap['C'] = map[byte]struct{}{'C': {}}
	patternCharMap['G'] = map[byte]struct{}{'G': {}}
	patternCharMap['T'] = map[byte]struct{}{'T': {}}

	// Degenerate bases
	patternCharMap['R'] = map[byte]struct{}{'A': {}, 'G': {}}  // puRine
	patternCharMap['Y'] = map[byte]struct{}{'C': {}, 'T': {}}  // pYrimidine
	patternCharMap['M'] = map[byte]struct{}{'A': {}, 'C': {}}  // aMino
	patternCharMap['K'] = map[byte]struct{}{'G': {}, 'T': {}}  // Keto
	patternCharMap['W'] = map[byte]struct{}{'A': {}, 'T': {}}  // Weak
	patternCharMap['S'] = map[byte]struct{}{'C': {}, 'G': {}}  // Strong
	patternCharMap['B'] = map[byte]struct{}{'C': {}, 'G': {}, 'T': {}}  // not A
	patternCharMap['D'] = map[byte]struct{}{'A': {}, 'G': {}, 'T': {}}  // not C
	patternCharMap['H'] = map[byte]struct{}{'A': {}, 'C': {}, 'T': {}}  // not G
	patternCharMap['V'] = map[byte]struct{}{'A': {}, 'C': {}, 'G': {}}  // not T
	patternCharMap['N'] = map[byte]struct{}{'A': {}, 'C': {}, 'G': {}, 'T': {}}  // aNy
}

// DNASearchEngine handles multiple DNA sequences with n-gram indexing
type DNASearchEngine struct {
	sequences  []string
	ngramIndex map[string]map[string]struct{} // ngram -> set of sequences
	n          int                             // n-gram size
}

// NewDNASearchEngine creates and preprocesses a search engine
func NewDNASearchEngine(sequences []string, n int) *DNASearchEngine {
	engine := &DNASearchEngine{
		sequences:  sequences,
		ngramIndex: make(map[string]map[string]struct{}),
		n:          n,
	}
	engine.buildIndex()
	return engine
}

// buildIndex preprocesses all sequences and builds n-gram index
func (e *DNASearchEngine) buildIndex() {
	fmt.Println("=== Building N-Gram Index ===")
	fmt.Printf("N-gram size: %d\n", e.n)
	fmt.Printf("Sequences to index: %d\n\n", len(e.sequences))

	for _, seq := range e.sequences {
		ngrams := e.extractNGrams(seq)
		fmt.Printf("Sequence: %s\n", seq)
		fmt.Printf("  N-grams: %v\n", ngrams)

		for _, ngram := range ngrams {
			if e.ngramIndex[ngram] == nil {
				e.ngramIndex[ngram] = make(map[string]struct{})
			}
			e.ngramIndex[ngram][seq] = struct{}{}
		}
	}

	fmt.Println("\n=== N-Gram Index ===")
	for ngram, seqs := range e.ngramIndex {
		seqList := []string{}
		for seq := range seqs {
			seqList = append(seqList, seq)
		}
		fmt.Printf("%s -> %v\n", ngram, seqList)
	}
	fmt.Println()
}

// extractNGrams returns all n-grams from a sequence
func (e *DNASearchEngine) extractNGrams(sequence string) []string {
	if len(sequence) < e.n {
		return []string{}
	}

	ngrams := []string{}
	for i := 0; i <= len(sequence)-e.n; i++ {
		ngrams = append(ngrams, sequence[i:i+e.n])
	}
	return ngrams
}

// expandWildcardNGram expands an n-gram with wildcards into all possible concrete n-grams
// This is the KEY function that leverages our original wildcard algorithm!
func (e *DNASearchEngine) expandWildcardNGram(ngram string) []string {
	// Check if ngram has any wildcards
	hasWildcard := false
	for i := 0; i < len(ngram); i++ {
		if ngram[i] != 'A' && ngram[i] != 'C' && ngram[i] != 'G' && ngram[i] != 'T' {
			hasWildcard = true
			break
		}
	}

	if !hasWildcard {
		return []string{ngram}
	}

	// Recursively expand wildcards
	return e.expandWildcardHelper(ngram, 0, "")
}

// expandWildcardHelper recursively generates all possible concrete n-grams
func (e *DNASearchEngine) expandWildcardHelper(ngram string, pos int, current string) []string {
	if pos == len(ngram) {
		return []string{current}
	}

	char := ngram[pos]
	possibleBases := patternCharMap[char]

	results := []string{}
	for base := range possibleBases {
		// Recursively expand rest of string
		subResults := e.expandWildcardHelper(ngram, pos+1, current+string(base))
		results = append(results, subResults...)
	}

	return results
}

// Search finds all sequences matching the query (WITH WILDCARD EXPANSION)
func (e *DNASearchEngine) Search(query string) []string {
	fmt.Printf("\n=== Searching for Query: %s ===\n", query)

	// Step 1: Extract n-grams from query
	queryNGrams := e.extractNGrams(query)
	fmt.Printf("Query n-grams: %v\n", queryNGrams)

	if len(queryNGrams) == 0 {
		// Query is shorter than n-gram size, fall back to direct search
		return e.directSearch(query)
	}

	// Step 1.5: Expand wildcards in n-grams
	fmt.Println("\nExpanding wildcards:")
	expandedNGrams := [][]string{}
	for _, ngram := range queryNGrams {
		expanded := e.expandWildcardNGram(ngram)
		fmt.Printf("  %s -> %v\n", ngram, expanded)
		expandedNGrams = append(expandedNGrams, expanded)
	}

	// Step 2: Look up each expanded n-gram and collect candidate sequences
	fmt.Println("\nN-gram lookups:")
	var candidates map[string]struct{}

	for i, ngramVariants := range expandedNGrams {
		// Union all sequences from all variants of this n-gram
		variantSeqs := make(map[string]struct{})
		for _, ngram := range ngramVariants {
			seqs := e.ngramIndex[ngram]
			for seq := range seqs {
				variantSeqs[seq] = struct{}{}
			}
		}

		seqList := []string{}
		for seq := range variantSeqs {
			seqList = append(seqList, seq)
		}
		fmt.Printf("  %v -> %v\n", ngramVariants, seqList)

		if i == 0 {
			// Initialize with first n-gram's sequences
			candidates = variantSeqs
		} else {
			// Intersect with current candidates
			newCandidates := make(map[string]struct{})
			for seq := range candidates {
				if _, exists := variantSeqs[seq]; exists {
					newCandidates[seq] = struct{}{}
				}
			}
			candidates = newCandidates
		}
	}

	// Step 3: Get candidate list
	candidateList := []string{}
	for seq := range candidates {
		candidateList = append(candidateList, seq)
	}
	fmt.Printf("\nCandidates after intersection: %v\n", candidateList)

	// Step 4: Filter false positives - verify full match
	fmt.Println("\nVerifying candidates:")
	results := []string{}
	for _, seq := range candidateList {
		if e.matchesQuery(seq, query) {
			fmt.Printf("  %s: ✓ MATCH\n", seq)
			results = append(results, seq)
		} else {
			fmt.Printf("  %s: ✗ FALSE POSITIVE\n", seq)
		}
	}

	return results
}

// directSearch performs direct substring matching for short queries
func (e *DNASearchEngine) directSearch(query string) []string {
	results := []string{}
	for _, seq := range e.sequences {
		if e.matchesQuery(seq, query) {
			results = append(results, seq)
		}
	}
	return results
}

// matchesQuery checks if sequence contains query (with wildcard support)
func (e *DNASearchEngine) matchesQuery(sequence string, query string) bool {
	if len(query) > len(sequence) {
		return false
	}

	// Sliding window search with wildcard matching
	for i := 0; i <= len(sequence)-len(query); i++ {
		matched := true
		for j := 0; j < len(query); j++ {
			seqChar := sequence[i+j]
			queryChar := query[j]

			// Check if query character matches sequence character (with wildcards)
			if _, ok := patternCharMap[queryChar][seqChar]; !ok {
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

func main() {
	initializePatternMap()

	// Upload DNA sequences
	sequences := []string{"GATTACA", "GATTG"}

	// Create search engine with n-gram size 4
	engine := NewDNASearchEngine(sequences, 4)

	// Test cases
	testCases := []struct {
		query       string
		description string
	}{
		{"GATT", "Exact match at beginning"},
		{"ATTACA", "Exact match in middle/end"},
		{"GATTR", "Wildcard R (A or G) - NOW FIXED!"},
		{"GATTM", "Wildcard M (A or C) - NOW FIXED!"},
		{"GATTRR", "Double wildcard RR (no match expected)"},
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("RUNNING SEARCH TESTS (WITH WILDCARD EXPANSION)")
	fmt.Println(strings.Repeat("=", 70))

	for _, tc := range testCases {
		fmt.Printf("\n%s\n", strings.Repeat("-", 70))
		fmt.Printf("Test: %s\n", tc.description)
		results := engine.Search(tc.query)
		fmt.Printf("\nRESULT: Matching sequences = %v\n", results)
	}

	// Demonstrate false positive example
	fmt.Printf("\n%s\n", strings.Repeat("=", 70))
	fmt.Println("FALSE POSITIVE EXAMPLE")
	fmt.Println(strings.Repeat("=", 70))

	falsePositiveSeqs := []string{"ATTAGATT"}
	fpEngine := NewDNASearchEngine(falsePositiveSeqs, 4)
	fpResults := fpEngine.Search("GATTA")
	fmt.Printf("\nRESULT: Matching sequences = %v\n", fpResults)
	fmt.Println("(Should be empty - GATTA is not in ATTAGATT as contiguous substring)")
}
