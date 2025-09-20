package main

func removeDuplicates(nums []int) int {
	duplicate := 0
	originalLen := len(nums)
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] == nums[i+1] {
			nums = append(nums[:i], nums[i+1:]...)
			duplicate++
			i--
		}
	}
	return originalLen - duplicate
}
