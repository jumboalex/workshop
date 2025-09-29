package tree

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func Flatten(root *TreeNode) {
	preOrder(root)
}

func preOrder(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}

	left := preOrder(root.Left)
	right := preOrder(root.Right)

	if left != nil {
		left.Right = root.Right
		root.Right = root.Left
		root.Left = nil
	}

	if right != nil {
		return right
	}
	if left != nil {
		return left
	}
	return root
}

func IsValidBST(root *TreeNode) bool {
	return validate(root, nil, nil)
}

func validate(root *TreeNode, low *int, high *int) bool {
	if root == nil {
		return true
	}
	if low != nil && root.Val <= *low {
		return false
	}
	if high != nil && root.Val >= *high {
		return false
	}
	return validate(root.Left, low, &root.Val) && validate(root.Right, &root.Val, high)
}
