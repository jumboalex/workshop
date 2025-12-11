package main

import (
	"fmt"
	"strings"
)

// Demonstration of different wildcard strategies
// (Simplified for illustration - not production code)

var wildcardMap = map[byte][]byte{
	'R': {'A', 'G'},
	'M': {'A', 'C'},
	'N': {'A', 'C', 'G', 'T'},
}

// ============================================================================
// STRATEGY 1: Wildcard Expansion (Current Implementation)
// ============================================================================
func strategy1_Expansion(pattern string) []string {
	fmt.Println("\n=== Strategy 1: Wildcard Expansion ===")
	fmt.Printf("Pattern: %s\n", pattern)

	variants := []string{""}
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		bases, isWildcard := wildcardMap[char]

		if isWildcard {
			newVariants := []string{}
			for _, variant := range variants {
				for _, base := range bases {
					newVariants = append(newVariants, variant+string(base))
				}
			}
			variants = newVariants
		} else {
			for j := range variants {
				variants[j] += string(char)
			}
		}
	}

	fmt.Printf("Expanded to %d variants: %v\n", len(variants), variants)
	return variants
}

// ============================================================================
// STRATEGY 2: Bit Vector (Conceptual)
// ============================================================================
func strategy2_BitVector(pattern string, sequence string) bool {
	fmt.Println("\n=== Strategy 2: Bit Vector ===")
	fmt.Printf("Pattern: %s, Sequence: %s\n", pattern, sequence)

	// Encode bases: A=00, C=01, G=10, T=11
	encode := func(base byte) uint8 {
		switch base {
		case 'A':
			return 0
		case 'C':
			return 1
		case 'G':
			return 2
		case 'T':
			return 3
		}
		return 0
	}

	// For wildcards, create bitmask
	for i := 0; i < len(pattern) && i < len(sequence); i++ {
		seqBits := encode(sequence[i])
		patBases, isWildcard := wildcardMap[pattern[i]]

		if isWildcard {
			match := false
			for _, base := range patBases {
				if encode(base) == seqBits {
					match = true
					break
				}
			}
			if !match {
				fmt.Println("Result: No match (bitwise comparison)")
				return false
			}
		} else {
			if encode(pattern[i]) != seqBits {
				fmt.Println("Result: No match (bitwise comparison)")
				return false
			}
		}
	}

	fmt.Println("Result: Match! (bitwise comparison)")
	return true
}

// ============================================================================
// STRATEGY 3: Trie-based Matching
// ============================================================================
type TrieNode struct {
	children  map[byte]*TrieNode
	sequences []string
}

func strategy3_Trie() {
	fmt.Println("\n=== Strategy 3: Trie with Wildcard Traversal ===")

	// Build simple trie
	root := &TrieNode{children: make(map[byte]*TrieNode)}

	// Insert "GATT"
	node := root
	for _, char := range "GATT" {
		if node.children[byte(char)] == nil {
			node.children[byte(char)] = &TrieNode{
				children: make(map[byte]*TrieNode),
			}
		}
		node = node.children[byte(char)]
	}
	node.sequences = []string{"GATT-sequence"}

	// Search for "GATR" (R = A or G)
	fmt.Println("Searching for 'GATR' (R = A or G)")
	fmt.Println("Traverse: G → A → T → [try both R branches: A and G]")

	// Navigate to "GAT"
	node = root
	for _, char := range "GAT" {
		node = node.children[byte(char)]
		if node == nil {
			fmt.Println("No match")
			return
		}
	}

	// At wildcard R, try both A and G
	results := []string{}
	for _, base := range []byte{'A', 'G'} {
		if child := node.children[base]; child != nil {
			results = append(results, child.sequences...)
		}
	}

	fmt.Printf("Found: %v\n", results)
}

// ============================================================================
// STRATEGY 4: Seed-and-Extend (BLAST-like)
// ============================================================================
func strategy4_SeedExtend(query string, sequences []string) []string {
	fmt.Println("\n=== Strategy 4: Seed-and-Extend ===")
	fmt.Printf("Query: %s\n", query)

	// Find longest exact substring (no wildcards)
	seed := ""
	for i := 0; i < len(query); i++ {
		for j := i + 1; j <= len(query); j++ {
			substring := query[i:j]
			hasWildcard := false
			for k := 0; k < len(substring); k++ {
				if _, isWild := wildcardMap[substring[k]]; isWild {
					hasWildcard = true
					break
				}
			}
			if !hasWildcard && len(substring) > len(seed) {
				seed = substring
			}
		}
	}

	fmt.Printf("Seed (exact substring): %s\n", seed)

	// Find candidates containing seed
	candidates := []string{}
	for _, seq := range sequences {
		if strings.Contains(seq, seed) {
			candidates = append(candidates, seq)
		}
	}

	fmt.Printf("Candidates from seed: %v\n", candidates)

	// Extend and verify with wildcard matching
	results := []string{}
	for _, seq := range candidates {
		// Find seed position and extend
		idx := strings.Index(seq, seed)
		if idx != -1 {
			// Verify full match (simplified)
			results = append(results, seq)
		}
	}

	fmt.Printf("After extension: %v\n", results)
	return results
}

// ============================================================================
// STRATEGY 5: Automaton (Conceptual)
// ============================================================================
func strategy5_Automaton() {
	fmt.Println("\n=== Strategy 5: Automaton (DFA/NFA) ===")
	fmt.Println("Pattern: ATTR (R = A or G)")
	fmt.Println("\nState Machine:")
	fmt.Println("  State 0 --A--> State 1")
	fmt.Println("  State 1 --T--> State 2")
	fmt.Println("  State 2 --T--> State 3")
	fmt.Println("  State 3 --A--> Accept")
	fmt.Println("  State 3 --G--> Accept")
	fmt.Println("\nMatching 'ATTA':")
	fmt.Println("  Start → A (State 1) → T (State 2) → T (State 3) → A (Accept) ✓")
	fmt.Println("\nMatching 'ATTG':")
	fmt.Println("  Start → A (State 1) → T (State 2) → T (State 3) → G (Accept) ✓")
}

// ============================================================================
// MAIN: Compare all strategies
// ============================================================================
func main() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("WILDCARD SEARCH STRATEGIES COMPARISON")
	fmt.Println(strings.Repeat("=", 70))

	// Test data
	sequences := []string{"GATTA", "GATTG", "GATTACA"}

	// Strategy 1: Expansion
	_ = strategy1_Expansion("ATTR")
	fmt.Println("Pros: Simple, uses existing index")
	fmt.Println("Cons: Exponential growth (NNNN = 256 variants!)")

	// Strategy 2: Bit Vector
	strategy2_BitVector("ATTR", "ATTA")
	fmt.Println("Pros: Very fast (CPU instructions)")
	fmt.Println("Cons: Limited to short patterns")

	// Strategy 3: Trie
	strategy3_Trie()
	fmt.Println("Pros: Natural wildcard traversal")
	fmt.Println("Cons: More memory, complex to build")

	// Strategy 4: Seed-and-Extend
	strategy4_SeedExtend("GATTRACA", sequences)
	fmt.Println("Pros: Proven by BLAST, good balance")
	fmt.Println("Cons: Needs good seed selection")

	// Strategy 5: Automaton
	strategy5_Automaton()
	fmt.Println("Pros: Efficient state machine")
	fmt.Println("Cons: Complex to build")

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("\n1. Expansion: Best for 1-2 wildcards (CURRENT)")
	fmt.Println("2. Bit Vector: Best for k-mer matching")
	fmt.Println("3. Trie: Best for prefix searches")
	fmt.Println("4. Seed-Extend: Best for real-world DNA (RECOMMENDED UPGRADE)")
	fmt.Println("5. Automaton: Best for complex patterns")
	fmt.Println("\nCurrent choice (Expansion) is GOOD for your use case! ✓")
}
