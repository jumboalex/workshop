package tree

import "testing"

func TestIsValidBST(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{
			"valid BST",
			&TreeNode{
				Val:   2,
				Left:  &TreeNode{Val: 1},
				Right: &TreeNode{Val: 3},
			},
			true,
		},
		{
			"invalid BST",
			&TreeNode{
				Val:  5,
				Left: &TreeNode{Val: 1},
				Right: &TreeNode{
					Val:   4,
					Left:  &TreeNode{Val: 3},
					Right: &TreeNode{Val: 6},
				},
			},
			false,
		},
		{
			"single node",
			&TreeNode{Val: 1},
			true,
		},
		{
			"nil tree",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidBST(tt.root)
			if got != tt.want {
				t.Errorf("IsValidBST() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("simple tree", func(t *testing.T) {
		root := &TreeNode{
			Val: 1,
			Left: &TreeNode{
				Val:   2,
				Left:  &TreeNode{Val: 3},
				Right: &TreeNode{Val: 4},
			},
			Right: &TreeNode{
				Val:   5,
				Right: &TreeNode{Val: 6},
			},
		}

		Flatten(root)

		// Verify it's flattened (all nodes should be on the right)
		expected := []int{1, 2, 3, 4, 5, 6}
		curr := root
		for i := 0; i < len(expected); i++ {
			if curr == nil {
				t.Fatalf("Tree ended early at index %d", i)
			}
			if curr.Val != expected[i] {
				t.Errorf("Node %d: got val %d, want %d", i, curr.Val, expected[i])
			}
			if curr.Left != nil {
				t.Errorf("Node %d has non-nil left child", i)
			}
			curr = curr.Right
		}
	})

	t.Run("nil tree", func(t *testing.T) {
		Flatten(nil) // Should not panic
	})
}
