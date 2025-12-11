package main

import (
	"container/heap"
	"fmt"
)

func main() {
	fmt.Println(lastStoneWeight([]int{2, 7, 4, 1, 8, 1})) // Expected: 1
	fmt.Println(lastStoneWeight([]int{1}))                 // Expected: 1
	fmt.Println(lastStoneWeight([]int{2, 2}))              // Expected: 0
}

// MaxHeap implements heap.Interface for a max heap of integers
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // Max heap: reverse comparison
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func lastStoneWeight(stones []int) int {
	// Build max heap from stones
	h := MaxHeap(stones)
	heap.Init(&h)

	// Keep smashing stones until at most one remains
	for h.Len() > 1 {
		// Get two heaviest stones
		y := heap.Pop(&h).(int) // Heaviest
		x := heap.Pop(&h).(int) // Second heaviest

		// If they're not equal, push the difference back
		if x != y {
			heap.Push(&h, y-x)
		}
		// If x == y, both are destroyed (nothing to push back)
	}

	// Return last stone weight, or 0 if no stones left
	if h.Len() == 0 {
		return 0
	}
	return h[0]
}
