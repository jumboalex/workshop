package main

import (
	"fmt"
	"os"

	arrayproblems "github.com/jumbo/workshop/practice/algorithms/array"
	"github.com/jumbo/workshop/practice/algorithms/linkedlist"
	stringproblems "github.com/jumbo/workshop/practice/algorithms/string"
	"github.com/jumbo/workshop/practice/algorithms/tree"
)

func main() {
	fmt.Println("=== Practice Problems Runner ===")
	fmt.Println()

	if len(os.Args) > 1 && os.Args[1] == "--help" {
		printHelp()
		return
	}

	removeDuplicates("abbaca")
	fmt.Println(removeDuplicatesStack("abbaca"))
	fmt.Println(minRemoveToMakeValid("lee(t(c)o)de)"))

	runArrayExamples()
	runStringExamples()
	runLinkedListExamples()
	runTreeExamples()
}

func removeDuplicates(s string) string {
	i := 0
	j := 1
	b := []byte(s)
	for j < len(b) {
		if b[i] == b[j] {
			b = append(b[:i], b[j+1:]...)
			if i > 0 {
				i--
				j--
			}
			continue
		}
		i++
		j++
	}
	return string(b)
}

func removeDuplicatesStack(s string) string {
	stack := []byte{}
	top := 0
	stack = append(stack, s[0])
	i := 1
	for i < len(s) {
		fmt.Println(string(stack[:top+1]), i)
		if s[i] == stack[top] {
			i++
			top--
			if top < 0 {
				top = 0
				stack[top] = s[i]
				i++
			}
		} else {
			top++
			if top >= len(stack) {
				stack = append(stack, s[i])
			} else {
				stack[top] = s[i]
			}
			i++

		}
	}
	return string(stack[:top+1])
}

func minRemoveToMakeValid(s string) string {
	lefts := []int{}
	top := -1
	rights := make(map[int]int)
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			top++
			lefts = append(lefts, i)
		}
		if s[i] == ')' {
			if top >= 0 {
				lefts = lefts[:top]
				top--
			} else {
				rights[i]++
			}
		}
	}
	for _, v := range lefts {
		rights[v]++
	}

	resultBytes := []byte{}
	for i := 0; i < len(s); i++ {
		if _, ok := rights[i]; ok {
			continue
		}
		resultBytes = append(resultBytes, s[i])
	}
	return string(resultBytes)
}

func printHelp() {
	fmt.Println("Usage: runner [--help]")
	fmt.Println("\nThis program demonstrates various algorithm implementations.")
	fmt.Println("Run without arguments to see all examples.")
}

func runArrayExamples() {
	fmt.Println("--- Array Problems ---")

	// PlusOne
	digits := []int{1, 2, 9}
	result := arrayproblems.PlusOne(digits)
	fmt.Printf("PlusOne(%v) = %v\n", []int{1, 2, 9}, result)

	// RemoveDuplicates
	nums := []int{1, 1, 2, 2, 3}
	length := arrayproblems.RemoveDuplicates(nums)
	fmt.Printf("RemoveDuplicates(%v) = %d (result: %v)\n", []int{1, 1, 2, 2, 3}, length, nums[:length])

	// CanPlaceFlowers
	flowerbed := []int{1, 0, 0, 0, 1}
	canPlace := arrayproblems.CanPlaceFlowers(flowerbed, 1)
	fmt.Printf("CanPlaceFlowers([1,0,0,0,1], 1) = %v\n", canPlace)

	// KidsWithCandies
	candies := []int{2, 3, 5, 1, 3}
	extraCandies := 3
	kidsResult := arrayproblems.KidsWithCandies(candies, extraCandies)
	fmt.Printf("KidsWithCandies(%v, %d) = %v\n", candies, extraCandies, kidsResult)

	// PrintDiagonalOrder
	fmt.Print("PrintDiagonalOrder([[1,2,3],[4,5,6],[7,8,9]]) = ")
	arrayproblems.PrintDiagonalOrder([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})
	arrayproblems.PrintDiagonalOrder([][]int{{1, 2}, {4, 5}, {7, 8}})
	arrayproblems.PrintDiagonalOrder([][]int{{1, 2, 3, 4}, {5, 6, 7, 8}})

	fmt.Println()
}

func runStringExamples() {
	fmt.Println("--- String Problems ---")

	// GcdOfStrings
	gcd := stringproblems.GcdOfStrings("ABCABC", "ABC")
	fmt.Printf("GcdOfStrings(\"ABCABC\", \"ABC\") = %q\n", gcd)

	// MergeAlternately
	merged := stringproblems.MergeAlternately("abc", "pqr")
	fmt.Printf("MergeAlternately(\"abc\", \"pqr\") = %q\n", merged)

	// AddBinary
	binarySum := stringproblems.AddBinary("1010", "1011")
	fmt.Printf("AddBinary(\"1010\", \"1011\") = %q\n", binarySum)

	// MaxVowels
	maxVow := stringproblems.MaxVowels("abciiidef", 3)
	fmt.Printf("MaxVowels(\"abciiidef\", 3) = %d\n", maxVow)

	fmt.Println()
}

func runLinkedListExamples() {
	fmt.Println("--- Linked List Problems ---")

	// Create a simple linked list: 1 -> 2 -> 3 -> 4
	head := &linkedlist.ListNode{Val: 1}
	head.Next = &linkedlist.ListNode{Val: 2}
	head.Next.Next = &linkedlist.ListNode{Val: 3}
	head.Next.Next.Next = &linkedlist.ListNode{Val: 4}

	fmt.Print("Original list: ")
	printList(head)

	linkedlist.ReorderList(head)
	fmt.Print("After ReorderList: ")
	printList(head)

	fmt.Println()
}

func runTreeExamples() {
	fmt.Println("--- Tree Problems ---")

	// Create a valid BST
	bst := &tree.TreeNode{
		Val:   2,
		Left:  &tree.TreeNode{Val: 1},
		Right: &tree.TreeNode{Val: 3},
	}

	isValid := tree.IsValidBST(bst)
	fmt.Printf("IsValidBST(tree with root=2, left=1, right=3) = %v\n", isValid)

	// Create a tree to flatten
	flattenTree := &tree.TreeNode{
		Val: 1,
		Left: &tree.TreeNode{
			Val:   2,
			Left:  &tree.TreeNode{Val: 3},
			Right: &tree.TreeNode{Val: 4},
		},
		Right: &tree.TreeNode{
			Val:   5,
			Right: &tree.TreeNode{Val: 6},
		},
	}

	fmt.Println("Flattening tree [1,2,5,3,4,null,6]...")
	tree.Flatten(flattenTree)
	fmt.Print("Flattened tree (preorder): ")
	printFlattenedTree(flattenTree)

	fmt.Println()
}

func printList(head *linkedlist.ListNode) {
	curr := head
	for curr != nil {
		fmt.Printf("%d", curr.Val)
		if curr.Next != nil {
			fmt.Print(" -> ")
		}
		curr = curr.Next
	}
	fmt.Println()
}

func printFlattenedTree(root *tree.TreeNode) {
	curr := root
	for curr != nil {
		fmt.Printf("%d", curr.Val)
		if curr.Right != nil {
			fmt.Print(" -> ")
		}
		curr = curr.Right
	}
	fmt.Println()
}
