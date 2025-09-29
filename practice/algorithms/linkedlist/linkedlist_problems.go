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
