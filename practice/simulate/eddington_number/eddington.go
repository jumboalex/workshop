package main

import (
	"fmt"
	"sort"
)

// GetEddingtonNumber calculates the Eddington number from an array of integers
// The Eddington number E is the largest number such that you have at least E values >= E
//
// Algorithm:
// 1. Sort the array in descending order
// 2. For each position i (0-indexed), check if the value at position i is >= i+1
// 3. The largest such i+1 is the Eddington number
//
// Example:
// - Array: [5, 4, 3, 2, 1] -> E = 3 (we have 3 values >= 3: 5, 4, 3)
// - Array: [10, 8, 5, 3, 2] -> E = 4 (we have 4 values >= 4: 10, 8, 5)
func GetEddingtonNumber(values []int) int {
	if len(values) == 0 {
		return 0
	}

	// Sort in descending order
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] > sorted[j]
	})

	// Find the largest E where we have at least E values >= E
	eddington := 0
	for i, value := range sorted {
		// i+1 represents how many values we've seen (1-indexed count)
		// If value >= i+1, then we have at least i+1 values >= i+1
		if value >= i+1 {
			eddington = i + 1
		} else {
			// Since array is sorted descending, we can break early
			break
		}
	}

	return eddington
}

// GetEddingtonNumberOptimized is an alternative implementation
// that doesn't break early (for comparison)
func GetEddingtonNumberOptimized(values []int) int {
	if len(values) == 0 {
		return 0
	}

	// Sort in descending order
	sorted := make([]int, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] > sorted[j]
	})

	// The Eddington number is the maximum i+1 where sorted[i] >= i+1
	eddington := 0
	for i := 0; i < len(sorted); i++ {
		// We need at least i+1 values >= i+1
		// sorted[i] is the (i+1)th largest value
		if sorted[i] >= i+1 {
			eddington = i + 1
		}
	}

	return eddington
}

func main() {
	fmt.Println("=== Eddington Number Calculator ===\n")

	testCases := []struct {
		name   string
		values []int
	}{
		{
			name:   "Example 1: [5, 4, 3, 2, 1]",
			values: []int{5, 4, 3, 2, 1},
		},
		{
			name:   "Example 2: [10, 8, 5, 3, 2]",
			values: []int{10, 8, 5, 3, 2},
		},
		{
			name:   "Example 3: [1, 1, 1, 1, 1]",
			values: []int{1, 1, 1, 1, 1},
		},
		{
			name:   "Example 4: [100, 50, 25, 10, 5, 3, 2, 1]",
			values: []int{100, 50, 25, 10, 5, 3, 2, 1},
		},
		{
			name:   "Example 5: [3, 3, 3, 3]",
			values: []int{3, 3, 3, 3},
		},
		{
			name:   "Example 6: [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]",
			values: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:   "Example 7: [0, 0, 0, 0]",
			values: []int{0, 0, 0, 0},
		},
		{
			name:   "Example 8: Empty array",
			values: []int{},
		},
		{
			name:   "Example 9: Single element [5]",
			values: []int{5},
		},
		{
			name:   "Example 10: [4, 4, 4, 4, 4]",
			values: []int{4, 4, 4, 4, 4},
		},
		{
			name:   "Cycling example: Daily rides",
			values: []int{10, 12, 15, 8, 20, 18, 14, 22, 25, 30},
		},
	}

	for _, tc := range testCases {
		e := GetEddingtonNumber(tc.values)
		fmt.Printf("%s\n", tc.name)
		fmt.Printf("  Values: %v\n", tc.values)
		fmt.Printf("  Eddington Number: %d\n", e)

		// Verify the result
		sorted := make([]int, len(tc.values))
		copy(sorted, tc.values)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] > sorted[j]
		})

		if e > 0 && e <= len(sorted) {
			fmt.Printf("  Verification: We have %d values >= %d: %v\n", e, e, sorted[:e])
			if e < len(sorted) {
				fmt.Printf("                But only %d values >= %d (next would be %d)\n", e, e+1, sorted[e])
			}
		}
		fmt.Println()
	}

	// Test both implementations produce same results
	fmt.Println("=== Comparing Both Implementations ===\n")
	testValues := [][]int{
		{5, 4, 3, 2, 1},
		{10, 8, 5, 3, 2},
		{3, 3, 3, 3},
	}

	for _, vals := range testValues {
		e1 := GetEddingtonNumber(vals)
		e2 := GetEddingtonNumberOptimized(vals)
		match := "✓"
		if e1 != e2 {
			match = "✗"
		}
		fmt.Printf("%v: E=%d (optimized=%d) %s\n", vals, e1, e2, match)
	}
}
