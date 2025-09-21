package main

func longestOnes(nums []int, k int) int {
	start := 0
	end := 0
	zeroes := 0
	maxLen := 0
	for end < len(nums) {
		if nums[end] == 0 {
			zeroes++
		}
		if zeroes > k {
			if nums[start] == 0 {
				zeroes--
			}
			start++
		}
		if end-start+1 > maxLen {
			maxLen = end - start + 1
		}

		end++

	}
	return maxLen
}
