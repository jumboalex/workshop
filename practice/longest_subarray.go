package main

import "fmt"

func longestSubarray(nums []int) int {
	fmt.Println("longestSubarray called")
	maxLen := 0

	start := 0
	end := 0
	delete := 0
	zeroes := 0
	for end < len(nums) {
		if nums[end] == 0 {
			zeroes++
			delete++
			for delete > 1 {
				if nums[start] == 0 {
					delete--
				}
				start++
			}
		}
		if end-start+1-delete > maxLen {
			maxLen = end - start + 1 - delete
		}
		end++
	}
	if zeroes == 0 {
		maxLen--
	}
	return maxLen
}
