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

type queueItem struct {
	node   *TreeNode
	column int
}

func VerticalOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	q := []queueItem{}
	columnTable := make(map[int][]int)
	minColumn := 0
	maxColumn := 0

	q = append(q, queueItem{root, 0})

	for len(q) > 0 {
		item := q[0]
		q = q[1:]

		columnTable[item.column] = append(columnTable[item.column], item.node.Val)

		if item.node.Left != nil {
			leftColumn := item.column - 1
			q = append(q, queueItem{item.node.Left, leftColumn})
			if leftColumn < minColumn {
				minColumn = leftColumn
			}
		}
		if item.node.Right != nil {
			rightColumn := item.column + 1
			q = append(q, queueItem{item.node.Right, rightColumn})
			if rightColumn > maxColumn {
				maxColumn = rightColumn
			}
		}
	}
	return buildResult(columnTable, minColumn, maxColumn)
}

func buildResult(columnTable map[int][]int, minColumn int, maxColumn int) [][]int {
	result := make([][]int, 0)
	for i := minColumn; i <= maxColumn; i++ {
		if values, ok := columnTable[i]; ok {
			result = append(result, values)
		} else {
			result = append(result, []int{})
		}
	}
	return result
}
