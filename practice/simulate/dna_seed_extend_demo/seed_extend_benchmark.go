package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
SEED-AND-EXTEND PERFORMANCE TEST

Compares seed-and-extend vs naive search to demonstrate the speedup
from avoiding unnecessary position checks.

NOTE: This file includes simplified versions of the seed-extend functions
to allow standalone execution for benchmarking.
*/

// IUPAC nucleotide codes
var iupacMap = map[byte]string{
	'A': "A", 'C': "C", 'G': "G", 'T': "T",
	'R': "AG", 'Y': "CT", 'M': "AC", 'K': "GT",
	'W': "AT", 'S': "CG", 'B': "CGT", 'D': "AGT",
	'H': "ACT", 'V': "ACG", 'N': "ACGT",
}

// matches checks if pattern character matches sequence character
func matches(patternChar, seqChar byte) bool {
	possibleBases, exists := iupacMap[patternChar]
	if !exists {
		return false
	}
	return strings.ContainsRune(possibleBases, rune(seqChar))
}

// Simplified engine for benchmarking
type SimpleSeedEngine struct {
	sequences  map[string]string
	seedIndex  map[string][]struct{ seqID string; pos int }
	minSeedLen int
}

func NewSimpleSeedEngine(minSeedLen int) *SimpleSeedEngine {
	return &SimpleSeedEngine{
		sequences:  make(map[string]string),
		seedIndex:  make(map[string][]struct{ seqID string; pos int }),
		minSeedLen: minSeedLen,
	}
}

func (e *SimpleSeedEngine) AddSequence(id, sequence string) {
	sequence = strings.ToUpper(sequence)
	e.sequences[id] = sequence
	for i := 0; i <= len(sequence)-e.minSeedLen; i++ {
		seed := sequence[i : i+e.minSeedLen]
		e.seedIndex[seed] = append(e.seedIndex[seed], struct{ seqID string; pos int }{id, i})
	}
}

func (e *SimpleSeedEngine) Search(query string) map[string][]int {
	query = strings.ToUpper(query)
	results := make(map[string][]int)

	// Extract seed (first contiguous non-wildcard segment)
	seed := ""
	seedOffset := 0
	for i := 0; i < len(query); i++ {
		char := query[i]
		if char == 'A' || char == 'C' || char == 'G' || char == 'T' {
			if seed == "" {
				seedOffset = i
			}
			seed += string(char)
			if len(seed) >= e.minSeedLen {
				break
			}
		} else if seed != "" {
			break
		}
	}

	if len(seed) < e.minSeedLen {
		// No good seed, fall back to naive
		return naiveSearch(e.sequences, query)
	}

	// Look up seed
	candidates := make(map[string]map[int]bool)
	if hits, exists := e.seedIndex[seed]; exists {
		for _, hit := range hits {
			queryStart := hit.pos - seedOffset
			if candidates[hit.seqID] == nil {
				candidates[hit.seqID] = make(map[int]bool)
			}
			candidates[hit.seqID][queryStart] = true
		}
	}

	// Extend and verify
	for seqID, positions := range candidates {
		sequence := e.sequences[seqID]
		for queryStart := range positions {
			if queryStart < 0 || queryStart+len(query) > len(sequence) {
				continue
			}
			match := true
			for i := 0; i < len(query); i++ {
				if !matches(query[i], sequence[queryStart+i]) {
					match = false
					break
				}
			}
			if match {
				results[seqID] = append(results[seqID], queryStart)
			}
		}
	}

	return results
}

// generateRandomDNA generates a random DNA sequence
func generateRandomDNA(length int) string {
	bases := []byte{'A', 'C', 'G', 'T'}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = bases[rand.Intn(4)]
	}
	return string(result)
}

// naiveSearch performs exhaustive search at every position
func naiveSearch(sequences map[string]string, query string) map[string][]int {
	results := make(map[string][]int)

	for seqID, sequence := range sequences {
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

// TestCase represents a performance test scenario
type TestCase struct {
	Name        string
	NumSeqs     int
	SeqLength   int
	Query       string
	Description string
}

func runPerformanceTest(tc TestCase) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("TEST: %s\n", tc.Name)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Setup: %d sequences × %d bases = %d total bases\n",
		tc.NumSeqs, tc.SeqLength, tc.NumSeqs*tc.SeqLength)
	fmt.Printf("Query: %s (%s)\n", tc.Query, tc.Description)

	// Generate test data
	rand.Seed(42) // Fixed seed for reproducibility
	sequences := make(map[string]string)
	for i := 0; i < tc.NumSeqs; i++ {
		seqID := fmt.Sprintf("seq%d", i)
		sequences[seqID] = generateRandomDNA(tc.SeqLength)
	}

	// Test 1: Seed-and-Extend
	fmt.Println("\n--- Seed-and-Extend ---")
	engine := NewSimpleSeedEngine(3)

	startBuild := time.Now()
	for seqID, seq := range sequences {
		engine.AddSequence(seqID, seq)
	}
	buildTime := time.Since(startBuild)

	fmt.Printf("Index build time: %v\n", buildTime)
	fmt.Printf("Index size: %d unique seeds\n", len(engine.seedIndex))

	startQuery := time.Now()
	results1 := engine.Search(tc.Query)
	queryTime1 := time.Since(startQuery)

	totalMatches1 := 0
	for _, positions := range results1 {
		totalMatches1 += len(positions)
	}

	fmt.Printf("Query time: %v\n", queryTime1)
	fmt.Printf("Matches found: %d\n", totalMatches1)

	// Test 2: Naive Search
	fmt.Println("\n--- Naive Search ---")
	startQuery2 := time.Now()
	results2 := naiveSearch(sequences, tc.Query)
	queryTime2 := time.Since(startQuery2)

	totalMatches2 := 0
	for _, positions := range results2 {
		totalMatches2 += len(positions)
	}

	fmt.Printf("Query time: %v\n", queryTime2)
	fmt.Printf("Matches found: %d\n", totalMatches2)

	// Comparison
	fmt.Println("\n--- Comparison ---")
	if totalMatches1 == totalMatches2 {
		fmt.Printf("✓ Results match (%d matches)\n", totalMatches1)
	} else {
		fmt.Printf("✗ Results differ! Seed: %d, Naive: %d\n", totalMatches1, totalMatches2)
	}

	if queryTime2 > queryTime1 {
		speedup := float64(queryTime2) / float64(queryTime1)
		fmt.Printf("⚡ Seed-and-Extend is %.2fx faster\n", speedup)
	} else {
		slowdown := float64(queryTime1) / float64(queryTime2)
		fmt.Printf("⚠️  Naive is %.2fx faster (unexpected!)\n", slowdown)
	}

	fmt.Printf("Time saved: %v\n", queryTime2-queryTime1)
}

func main() {
	fmt.Println("SEED-AND-EXTEND PERFORMANCE BENCHMARKS")
	fmt.Println("Testing various scenarios to show algorithm behavior\n")

	testCases := []TestCase{
		{
			Name:        "Small Database, Specific Query",
			NumSeqs:     10,
			SeqLength:   1000,
			Query:       "ATTR",
			Description: "ATT seed is specific, few hits",
		},
		{
			Name:        "Medium Database, Specific Query",
			NumSeqs:     100,
			SeqLength:   10000,
			Query:       "GATTACA",
			Description: "Long seed, very specific",
		},
		{
			Name:        "Large Database, Common Query",
			NumSeqs:     50,
			SeqLength:   50000,
			Query:       "ATW", // W = A or T
			Description: "AT seed is common, many hits",
		},
		{
			Name:        "Small Database, Multiple Wildcards",
			NumSeqs:     20,
			SeqLength:   5000,
			Query:       "ARNTY", // R = A/G, N = any, Y = C/T
			Description: "AR seed moderately specific",
		},
		{
			Name:        "Worst Case - All Wildcards",
			NumSeqs:     10,
			SeqLength:   1000,
			Query:       "NNNN",
			Description: "No seeds, falls back to naive",
		},
	}

	for _, tc := range testCases {
		runPerformanceTest(tc)
	}

	// Summary
	fmt.Println("\n\n" + strings.Repeat("=", 70))
	fmt.Println("KEY TAKEAWAYS")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(`
1. SEED QUALITY MATTERS
   - Specific seeds (e.g., "GATTACA") = Few candidates = Fast
   - Common seeds (e.g., "AT") = Many candidates = Slower
   - No seeds (all wildcards) = Falls back to naive = Slowest

2. DATABASE SIZE
   - Larger databases amplify the speedup
   - Naive checks every position: O(N × M × Q)
   - Seed-extend checks only candidates: O(H × Q) where H << N×M

3. WHEN SEED-EXTEND WINS
   ✓ Queries with at least one specific seed (3+ exact bases)
   ✓ Large databases (where avoiding checks matters)
   ✓ Wildcards scattered in the pattern (not all at once)

4. WHEN SEED-EXTEND STRUGGLES
   ✗ All-wildcard queries (no seeds to anchor on)
   ✗ Very common seeds (e.g., "AA" in A/T-rich sequences)
   ✗ Tiny databases (naive is fast enough anyway)

5. REAL-WORLD USAGE
   - BLAST searches billions of bases in seconds using this approach
   - The key is the seed filtering step dramatically reduces candidates
   - For typical DNA queries, expect 10-100× speedup over naive search
`)

	fmt.Println("Try experimenting with different query patterns to see the effect!")
}
