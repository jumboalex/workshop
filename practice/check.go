package main

import (
	"fmt"
)

func main() {
	checkPlusOne()
	readFile()
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
}
