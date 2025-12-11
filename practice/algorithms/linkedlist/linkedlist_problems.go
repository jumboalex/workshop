package linkedlist

type ListNode struct {
	Val  int
	Next *ListNode
}

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func ReorderList(head *ListNode) {
	if head == nil || head.Next == nil {
		return
	}

	// Find middle
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Reverse second half
	var prev *ListNode
	curr := slow.Next
	slow.Next = nil

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	// Merge two halves
	first, second := head, prev
	for second != nil {
		tmp1, tmp2 := first.Next, second.Next
		first.Next = second
		second.Next = tmp1
		first = tmp1
		second = tmp2
	}
}

func CopyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}

	nodeList := []*Node{}
	for p := head; p != nil; p = p.Next {
		nodeList = append(nodeList, p)
	}

	randomMap := make(map[int]int)
	for i, node := range nodeList {
		r := node.Random
		if r == nil {
			randomMap[i] = -1
			continue
		}
		for j, n := range nodeList {
			if n == r {
				randomMap[i] = j
				break
			}
		}
	}

	newHead := &Node{Val: head.Val}
	newList := []*Node{newHead}
	for i := 1; i < len(nodeList); i++ {
		newNode := &Node{Val: nodeList[i].Val}
		newList[i-1].Next = newNode
		newList = append(newList, newNode)
	}
	for i, newNode := range newList {
		r := randomMap[i]
		if r == -1 {
			newNode.Random = nil
		} else {
			newNode.Random = newList[r]
		}
	}

	return newHead
}

func MergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	current := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			current.Next = list1
			list1 = list1.Next
		} else {
			current.Next = list2
			list2 = list2.Next
		}
		current = current.Next
	}

	if list1 != nil {
		current.Next = list1
	}
	if list2 != nil {
		current.Next = list2
	}

	return dummy.Next
}

func ReorderList2(head *ListNode) {
	if head == nil || head.Next == nil {
		return
	}

	count := 0
	for i := head; i != nil; i = i.Next {
		count++
	}
	i := head
	for step := 0; step < count/2; step++ {
		j := i
		for j.Next != nil && j.Next.Next != nil {
			j = j.Next
		}
		if j.Next == nil {
			break
		}
		k := j.Next
		j.Next = nil
		k.Next = i.Next
		i.Next = k
		i = k.Next
	}

}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	nodeL := []*ListNode{}
	node := head
	for node != nil {
		nodeL = append(nodeL, node)
		node = node.Next
	}

	// Edge case: empty list or invalid n
	if len(nodeL) == 0 || n <= 0 || n > len(nodeL) {
		return head
	}

	// Special case: removing the first node (head)
	if n == len(nodeL) {
		return head.Next
	}

	// General case: remove node at position (len-n) by updating predecessor
	nodeL[len(nodeL)-n-1].Next = nodeL[len(nodeL)-n].Next
	return head
}

func removeElements(head *ListNode, val int) *ListNode {
	if head == nil {
		return head
	}

	p := &ListNode{Next: head}
	q := p
	for q != nil && q.Next != nil {
		if q.Next.Val == val {
			q.Next = q.Next.Next
			// Don't move q forward - stay and check the new q.Next
		} else {
			q = q.Next // Only move forward if no deletion
		}
	}
	return p.Next
}

type NodeMultiLevel struct {
	Val   int
	Prev  *NodeMultiLevel
	Next  *NodeMultiLevel
	Child *NodeMultiLevel
}

func flatten(root *NodeMultiLevel) *NodeMultiLevel {
	if root == nil {
		return root
	}

	p := root
	for p != nil {
		// If current node has a child
		if p.Child != nil {
			next := p.Next

			// Recursively flatten the child list (handles nested children)
			child := flatten(p.Child)

			// Find the tail of the flattened child list
			tail := child
			for tail.Next != nil {
				tail = tail.Next
			}

			// Insert flattened child list between p and next
			p.Next = child
			child.Prev = p
			tail.Next = next
			if next != nil {
				next.Prev = tail
			}

			// Clear the child pointer (important!)
			p.Child = nil
		}
		p = p.Next
	}
	return root
}

type NodeCircular struct {
	Val  int
	Next *NodeCircular
}

func insert(aNode *NodeCircular, x int) *NodeCircular {
	n := &NodeCircular{Val: x}

	// Case 1: Empty list
	if aNode == nil {
		n.Next = n
		return n
	}

	// Case 2: Single node
	if aNode.Next == aNode {
		n.Next = aNode
		aNode.Next = n
		return aNode
	}

	// Case 3: Multiple nodes - find insertion point
	p := aNode
	for {
		// Case 3a: Insert between p and p.Next when value fits in range
		if p.Val <= x && x <= p.Next.Val {
			n.Next = p.Next
			p.Next = n
			return aNode
		}

		// Case 3b: Insert at the boundary (between max and min)
		// This happens when p.Val > p.Next.Val (we're at the wrap-around point)
		// AND (x is larger than max OR x is smaller than min)
		if p.Val > p.Next.Val {
			if x >= p.Val || x <= p.Next.Val {
				n.Next = p.Next
				p.Next = n
				return aNode
			}
		}

		p = p.Next

		// If we've completed a full circle, insert after current position
		// This handles the case where all values are equal
		if p == aNode {
			n.Next = p.Next
			p.Next = n
			return aNode
		}
	}
}

func rotateRight(head *ListNode, k int) *ListNode {
	// Edge case: empty list or single node
	if head == nil || head.Next == nil {
		return head
	}

	// Step 1: Find length and tail of the list
	length := 1
	tail := head
	for tail.Next != nil {
		tail = tail.Next
		length++
	}

	// Step 2: Optimize k (no need to rotate more than length times)
	k = k % length
	if k == 0 {
		return head // No rotation needed
	}

	// Step 3: Find the new tail (at position length - k - 1)
	newTailPos := length - k - 1
	newTail := head
	for i := 0; i < newTailPos; i++ {
		newTail = newTail.Next
	}

	// Step 4: The new head is after the new tail
	newHead := newTail.Next

	// Step 5: Break the old list and connect tail to old head
	newTail.Next = nil
	tail.Next = head

	return newHead
}
