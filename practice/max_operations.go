package main

import (
	"fmt"
	"sort"
)

func maxOperations(nums []int, k int) int {
	sort.Ints(nums)
	i := 0
	j := len(nums) - 1
	result := 0
	fmt.Println("nums:", nums, "i:", i, " j:", j)
	for i < j {
		sum := nums[i] + nums[j]
		if sum == k {
			result++
			i++
			j--
		} else if sum < k {
			i++
		} else {
			j--
		}
	}
	return result
}
