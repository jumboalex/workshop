package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func flatten(root *TreeNode) {
	preOrder(root)
}

func preOrder(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}
	fmt.Println("preOrder called with:", root.Val)
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
