package linkedlist

import (
	"reflect"
	"testing"
)

func TestReorderList(t *testing.T) {
	tests := []struct {
		name string
		vals []int
		want []int
	}{
		{"example 1", []int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{"example 2", []int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
		{"single", []int{1}, []int{1}},
		{"two elements", []int{1, 2}, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.vals)
			ReorderList(head)
			got := listToSlice(head)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReorderList(%v) = %v, want %v", tt.vals, got, tt.want)
			}
		})
	}
}

func TestCopyRandomList(t *testing.T) {
	t.Run("simple list", func(t *testing.T) {
		// Create: [[7,null],[13,0],[11,4],[10,2],[1,0]]
		node1 := &Node{Val: 7}
		node2 := &Node{Val: 13}
		node3 := &Node{Val: 11}
		node4 := &Node{Val: 10}
		node5 := &Node{Val: 1}

		node1.Next = node2
		node2.Next = node3
		node3.Next = node4
		node4.Next = node5

		node1.Random = nil
		node2.Random = node1
		node3.Random = node5
		node4.Random = node3
		node5.Random = node1

		copied := CopyRandomList(node1)

		// Verify it's a deep copy
		if copied == node1 {
			t.Error("CopyRandomList returned same node, expected deep copy")
		}

		// Verify values
		original := []int{7, 13, 11, 10, 1}
		curr := copied
		i := 0
		for curr != nil {
			if curr.Val != original[i] {
				t.Errorf("Node %d: got val %d, want %d", i, curr.Val, original[i])
			}
			curr = curr.Next
			i++
		}
	})

	t.Run("nil list", func(t *testing.T) {
		got := CopyRandomList(nil)
		if got != nil {
			t.Errorf("CopyRandomList(nil) = %v, want nil", got)
		}
	})
}

// Helper functions
func createList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	curr := head
	for i := 1; i < len(vals); i++ {
		curr.Next = &ListNode{Val: vals[i]}
		curr = curr.Next
	}
	return head
}

func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}
