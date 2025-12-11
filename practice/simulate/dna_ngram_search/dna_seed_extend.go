package main

import (
	"fmt"
	"strings"
)

/*
SEED-AND-EXTEND ALGORITHM FOR DNA WILDCARD SEARCH

This is the approach used by BLAST (Basic Local Alignment Search Tool) and
other bioinformatics tools for efficient sequence matching with wildcards.

ALGORITHM OVERVIEW:
==================

1. SEED PHASE: Find exact-match "seeds" (subsequences without wildcards)
   - Extract all non-wildcard segments from the query
   - Use n-gram index to quickly locate these seeds in sequences
   - Seeds act as anchor points for potential matches

2. EXTEND PHASE: Expand from seeds to verify full match
   - For each seed hit, extend left and right
   - Check if the full pattern (including wildcards) matches
   - This avoids checking every position in the sequence

EXAMPLE:
========
Query: "ATTR" where R = A or G
Seeds: "ATT" (3-mer without wildcards)

Sequence: "GATTACATTAGC"
           01234567890123

Step 1: Find "ATT" seeds
   - Found at position 1: "GATTACATTAGC"
                            ^ATT
   - Found at position 6: "GATTACATTAGC"
                                ^ATT

Step 2: Extend from each seed
   - Position 1: Check "ATTA" - matches ATTR? YES (A matches R)
   - Position 6: Check "ATTA" - matches ATTR? YES (A matches R)

Results: ["GATTACATTAGC"] with matches at positions [1, 6]

ADVANTAGES:
===========
1. Speed: Only checks positions near seeds (not every position)
2. Memory: Uses compact n-gram index for seeds only
3. Wildcards: Handles any number of wildcards naturally
4. Proven: Used by BLAST in production for decades

COMPLEXITY:
===========
- Build Index: O(N * M) where N = total chars, M = sequence length
- Query: O(S + H * Q) where:
  - S = seed lookup time (O(1) with hash table)
  - H = number of hits (seeds found)
  - Q = query length (extension check)
- Space: O(N * M) for n-gram index

WHEN TO USE:
============
✓ Queries with mix of exact bases and wildcards
✓ Long sequences (genomic data)
✓ Multiple wildcards scattered in pattern
✓ Need balance of speed and memory
✗ All wildcards (no seeds to anchor on)
*/

// IUPAC nucleotide codes for wildcards
var iupacMap = map[byte]string{
	'A': "A",
	'C': "C",
	'G': "G",
	'T': "T",
	'R': "AG",   // puRine
	'Y': "CT",   // pYrimidine
	'M': "AC",   // aMino
	'K': "GT",   // Keto
	'W': "AT",   // Weak
	'S': "CG",   // Strong
	'B': "CGT",  // not A
	'D': "AGT",  // not C
	'H': "ACT",  // not G
	'V': "ACG",  // not T
	'N': "ACGT", // aNy
}

// Seed represents a non-wildcard segment from the query
type Seed struct {
	Sequence string // The actual seed sequence (no wildcards)
	Offset   int    // Position in the original query
}

// SeedMatch represents where a seed was found
type SeedMatch struct {
	SequenceID string // Which sequence contained the match
	Position   int    // Position in the sequence
	SeedOffset int    // Which seed this came from
}

// DNASeedExtendEngine performs seed-and-extend searching
type DNASeedExtendEngine struct {
	sequences  map[string]string            // ID -> sequence
	seedIndex  map[string][]SeedMatch       // seed -> matches
	minSeedLen int                          // Minimum seed length
}

// NewDNASeedExtendEngine creates a new seed-and-extend search engine
func NewDNASeedExtendEngine(minSeedLen int) *DNASeedExtendEngine {
	return &DNASeedExtendEngine{
		sequences:  make(map[string]string),
		seedIndex:  make(map[string][]SeedMatch),
		minSeedLen: minSeedLen,
	}
}

// AddSequence adds a sequence to the database and indexes it
func (e *DNASeedExtendEngine) AddSequence(id string, sequence string) {
	sequence = strings.ToUpper(sequence)
	e.sequences[id] = sequence

	// Index all k-mers as potential seeds
	for i := 0; i <= len(sequence)-e.minSeedLen; i++ {
		seed := sequence[i : i+e.minSeedLen]
		e.seedIndex[seed] = append(e.seedIndex[seed], SeedMatch{
			SequenceID: id,
			Position:   i,
			SeedOffset: 0, // Will be set during search
		})
	}
}

// extractSeeds finds all non-wildcard segments in the query
func (e *DNASeedExtendEngine) extractSeeds(query string) []Seed {
	query = strings.ToUpper(query)
	var seeds []Seed

	start := -1
	for i := 0; i < len(query); i++ {
		char := query[i]
		isWildcard := char != 'A' && char != 'C' && char != 'G' && char != 'T'

		if !isWildcard {
			if start == -1 {
				start = i // Start of new seed
			}
		} else {
			if start != -1 {
				// End of seed
				seedSeq := query[start:i]
				if len(seedSeq) >= e.minSeedLen {
					seeds = append(seeds, Seed{
						Sequence: seedSeq,
						Offset:   start,
					})
				}
				start = -1
			}
		}
	}

	// Handle seed at end of query
	if start != -1 {
		seedSeq := query[start:]
		if len(seedSeq) >= e.minSeedLen {
			seeds = append(seeds, Seed{
				Sequence: seedSeq,
				Offset:   start,
			})
		}
	}

	return seeds
}

// matches checks if a pattern character matches a sequence character
func matches(patternChar, seqChar byte) bool {
	possibleBases, exists := iupacMap[patternChar]
	if !exists {
		return false
	}
	return strings.ContainsRune(possibleBases, rune(seqChar))
}

// extendMatch tries to match the full query at the given position
func (e *DNASeedExtendEngine) extendMatch(sequence, query string, seedPos, seedOffset int) bool {
	// Calculate where the query would start in the sequence
	queryStart := seedPos - seedOffset

	// Check if query fits in sequence
	if queryStart < 0 || queryStart+len(query) > len(sequence) {
		return false
	}

	// Check each position
	for i := 0; i < len(query); i++ {
		seqChar := sequence[queryStart+i]
		patternChar := query[i]

		if !matches(patternChar, seqChar) {
			return false
		}
	}

	return true
}

// Search finds all sequences containing the query using seed-and-extend
func (e *DNASeedExtendEngine) Search(query string) map[string][]int {
	query = strings.ToUpper(query)
	results := make(map[string][]int)

	// Step 1: Extract seeds from query
	seeds := e.extractSeeds(query)

	if len(seeds) == 0 {
		// No seeds found - query is all wildcards
		// Fall back to checking every position (rare case)
		fmt.Println("Warning: No seeds found in query, checking all positions")
		return e.searchWithoutSeeds(query)
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Extracted %d seed(s):\n", len(seeds))
	for _, seed := range seeds {
		fmt.Printf("  - '%s' at offset %d\n", seed.Sequence, seed.Offset)
	}

	// Step 2: Find seed matches
	candidateMatches := make(map[string]map[int]bool)

	for _, seed := range seeds {
		matches, exists := e.seedIndex[seed.Sequence]
		if !exists {
			continue
		}

		fmt.Printf("\nSeed '%s' found %d time(s):\n", seed.Sequence, len(matches))

		for _, match := range matches {
			if candidateMatches[match.SequenceID] == nil {
				candidateMatches[match.SequenceID] = make(map[int]bool)
			}

			// Calculate where the full query would start
			queryStart := match.Position - seed.Offset
			candidateMatches[match.SequenceID][queryStart] = true

			fmt.Printf("  - Sequence '%s' at position %d (query would start at %d)\n",
				match.SequenceID, match.Position, queryStart)
		}
	}

	// Step 3: Extend and verify each candidate
	fmt.Printf("\nExtending candidates:\n")
	for seqID, positions := range candidateMatches {
		sequence := e.sequences[seqID]

		for queryStart := range positions {
			if queryStart < 0 || queryStart+len(query) > len(sequence) {
				continue
			}

			// Try to match the full query at this position
			if e.extendMatch(sequence, query, queryStart+seeds[0].Offset, seeds[0].Offset) {
				results[seqID] = append(results[seqID], queryStart)
				matchedSubseq := sequence[queryStart : queryStart+len(query)]
				fmt.Printf("  ✓ Sequence '%s' at position %d: %s matches %s\n",
					seqID, queryStart, matchedSubseq, query)
			} else {
				fmt.Printf("  ✗ Sequence '%s' at position %d: no match on extension\n",
					seqID, queryStart)
			}
		}
	}

	return results
}

// searchWithoutSeeds is a fallback for queries with no seeds (all wildcards)
func (e *DNASeedExtendEngine) searchWithoutSeeds(query string) map[string][]int {
	results := make(map[string][]int)

	for seqID, sequence := range e.sequences {
		for i := 0; i <= len(sequence)-len(query); i++ {
			match := true
			for j := 0; j < len(query); j++ {
				if !matches(query[j], sequence[i+j]) {
					match = false
					break
				}
			}
			if match {
				results[seqID] = append(results[seqID], i)
			}
		}
	}

	return results
}

// PrintResults displays search results in a readable format
func PrintResults(results map[string][]int) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("SEARCH RESULTS")
	fmt.Println(strings.Repeat("=", 60))

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	for seqID, positions := range results {
		fmt.Printf("\nSequence: %s\n", seqID)
		fmt.Printf("Matches at positions: %v\n", positions)
	}
}

func main() {
	fmt.Println("SEED-AND-EXTEND DNA WILDCARD SEARCH DEMO")
	fmt.Println(strings.Repeat("=", 60))

	// Create search engine with minimum seed length of 3
	engine := NewDNASeedExtendEngine(3)

	// Add test sequences
	fmt.Println("\nAdding sequences to database...")
	engine.AddSequence("seq1", "GATTACATTAGC")
	engine.AddSequence("seq2", "CCGATTAGGATT")
	engine.AddSequence("seq3", "TTTTATTGCCCC")

	fmt.Printf("Indexed %d sequences\n", len(engine.sequences))
	fmt.Printf("Indexed %d unique seeds\n", len(engine.seedIndex))

	// Test cases
	testCases := []string{
		"ATTR",     // R = A or G
		"ATYN",     // Y = C or T, N = any
		"GATTR",    // Longer pattern
		"NNNN",     // All wildcards (worst case)
		"ATTW",     // W = A or T
	}

	for _, query := range testCases {
		fmt.Println("\n" + strings.Repeat("-", 60))
		results := engine.Search(query)
		PrintResults(results)
	}

	// Detailed example showing the process
	fmt.Println("\n\n" + strings.Repeat("=", 60))
	fmt.Println("DETAILED WALKTHROUGH: Query 'ATTR' in 'GATTACATTAGC'")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n1. SEED EXTRACTION:")
	fmt.Println("   Query: ATTR")
	fmt.Println("   A, T, T are exact bases (no wildcards)")
	fmt.Println("   R is a wildcard (A or G)")
	fmt.Println("   Seed extracted: 'ATT' at offset 0")

	fmt.Println("\n2. SEED LOOKUP:")
	fmt.Println("   Looking up 'ATT' in index...")
	fmt.Println("   Sequence: GATTACATTAGC")
	fmt.Println("            0123456789012")
	fmt.Println("             ^^^     ^^^")
	fmt.Println("   Found at positions: 1, 6")

	fmt.Println("\n3. EXTENSION:")
	fmt.Println("   Position 1:")
	fmt.Println("   Sequence: GATTACATTAGC")
	fmt.Println("             .ATTA......")
	fmt.Println("   Query:     ATTR")
	fmt.Println("   Check: A=A ✓, T=T ✓, T=T ✓, R=A ✓")
	fmt.Println("   Result: MATCH!")

	fmt.Println("\n   Position 6:")
	fmt.Println("   Sequence: GATTACATTAGC")
	fmt.Println("             ......ATTA..")
	fmt.Println("   Query:           ATTR")
	fmt.Println("   Check: A=A ✓, T=T ✓, T=T ✓, R=A ✓")
	fmt.Println("   Result: MATCH!")

	fmt.Println("\n4. FINAL RESULT:")
	fmt.Println("   Sequence 'seq1' has matches at positions [1, 6]")
}
