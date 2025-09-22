package main

func isValidBST(root *TreeNode) bool {
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
