package main

type MyLinkedList struct {
	Head *Node
}

type Node struct {
	Val  int
	Next *Node
	Prev *Node
}

func Constructor() MyLinkedList {
	return MyLinkedList{}
}

func (this *MyLinkedList) Get(index int) int {
	n := this.Head
	for i := 0; i < index; i++ {
		if n != nil {
			n = n.Next
		}
	}
	if n == nil {
		return -1
	}
	return n.Val
}

func (this *MyLinkedList) AddAtHead(val int) {
	node := &Node{Val: val}
	node.Next = this.Head
	if this.Head != nil {
		this.Head.Prev = node
	}
	this.Head = node
}

func (this *MyLinkedList) AddAtTail(val int) {
	node := &Node{Val: val}
	if this.Head == nil {
		this.Head = node
		return
	}
	n := this.Head
	for n != nil && n.Next != nil {
		n = n.Next
	}

	n.Next = node
	node.Prev = n
}

func (this *MyLinkedList) AddAtIndex(index int, val int) {
	if index == 0 {
		this.AddAtHead(val)
		return
	}

	node := &Node{Val: val}
	n := this.Head
	for i := 0; i < index-1; i++ {
		if n == nil {
			return
		}
		n = n.Next
	}
	if n == nil {
		return
	}

	node.Next = n.Next
	node.Prev = n
	n.Next = node

	// Update the Prev pointer of the next node (if it exists)
	if node.Next != nil {
		node.Next.Prev = node
	}
}

func (this *MyLinkedList) DeleteAtIndex(index int) {
	if this.Head == nil {
		return
	}
	if index == 0 {
		this.Head = this.Head.Next
		if this.Head != nil {
			this.Head.Prev = nil
		}
		return
	}
	n := this.Head
	for i := 0; i < index-1; i++ {
		if n == nil || n.Next == nil {
			return
		}
		n = n.Next
	}
	if n.Next == nil {
		return
	}

	n.Next = n.Next.Next
	// Update Prev pointer only if the new next node exists
	if n.Next != nil {
		n.Next.Prev = n
	}
}
