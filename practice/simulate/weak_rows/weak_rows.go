package main

import "container/heap"

type RowDim struct {
	Soldiers int
	Index    int
}

type MinHeap []RowDim

func (h MinHeap) Len() int {
	return len(h)
}
func (h MinHeap) Less(i, j int) bool {
	if h[i].Soldiers == h[j].Soldiers {
		return h[i].Index < h[j].Index
	} else {
		return h[i].Soldiers < h[j].Soldiers
	}
}
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(RowDim))
}
func (h *MinHeap) Pop() any {
	old := *h
	l := len(old)
	x := old[l-1]
	*h = old[0 : l-1]
	return x
}
func kWeakestRows(mat [][]int, k int) []int {
	rows := []RowDim{}
	for i, r := range mat {
		soldiers := 0
		rowDim := RowDim{}
		for _, p := range r {
			if p == 1 {
				soldiers++
			}
		}
		rowDim.Soldiers = soldiers
		rowDim.Index = i
		rows = append(rows, rowDim)
	}

	h := MinHeap(rows)
	heap.Init(&h)

	result := []int{}
	for i := 0; i < k; i++ {
		x := heap.Pop(&h).(RowDim)
		result = append(result, x.Index)
	}
	return result
}
