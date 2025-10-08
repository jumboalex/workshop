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
