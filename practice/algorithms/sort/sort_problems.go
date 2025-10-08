package sortproblems

import "container/heap"

func HeightChecker(heights []int) int {
	expected := make([]int, len(heights))
	copy(expected, heights)
	hasSwap := true
	for hasSwap {
		hasSwap = false

		for i := 0; i < len(expected)-1; i++ {
			if expected[i] > expected[i+1] {
				expected[i], expected[i+1] = expected[i+1], expected[i]
				hasSwap = true
			}
		}
	}
	result := 0
	for i := 0; i < len(expected); i++ {
		if expected[i] != heights[i] {
			result++
		}
	}
	return result
}

func sortArray(nums []int) []int {
	n := len(nums)

	// Build max heap
	for i := n/2 - 1; i >= 0; i-- {
		heapify(nums, n, i)
	}

	// Extract elements from heap one by one
	for i := n - 1; i > 0; i-- {
		nums[0], nums[i] = nums[i], nums[0]
		heapify(nums, i, 0)
	}

	return nums
}

func heapify(arr []int, n int, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	if right < n && arr[right] > arr[largest] {
		largest = right
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func SortArrayWithHeap(nums []int) []int {
	h := &IntHeap{}
	heap.Init(h)

	for _, num := range nums {
		heap.Push(h, num)
	}

	result := make([]int, 0, len(nums))
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(int))
	}

	return result
}
