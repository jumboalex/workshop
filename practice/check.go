package main

import (
	"fmt"
)

func main() {
	readFile()
	checkPlusOne()
	fmt.Println(MultiplyString("123", "456"))
	fmt.Println(addBinary("11", "1"))
	fmt.Println(removeDuplicates([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))
	fmt.Println(gcdOfStrings("ABABAB", "ABAB"))
	fmt.Println(canPlaceFlowers([]int{0, 0, 0, 0, 0, 1, 0, 0}, 1))
	fmt.Println(maxOperations([]int{3, 1, 3, 4, 3}, 6))
	fmt.Println(maxVowels("aeiou", 2))
	fmt.Println(longestSubarray([]int{0, 1, 1, 1, 0, 1, 1, 0, 1}))
}

func checkPlusOne() {
	var tests [][]int
	tests = append(tests, []int{1, 2, 3})
	tests = append(tests, []int{4, 3, 2, 1})
	tests = append(tests, []int{9})
	tests = append(tests, []int{9, 9, 9})

	for _, test := range tests {
		fmt.Println(PlusOne(test))
	}

	num1 := "123"
	for i := len(num1) - 1; i >= 0; i-- {
		fmt.Println(num1[i] - '0')
	}
}
