package arrayproblems

import "sort"

func CanPlaceFlowers(flowerbed []int, n int) bool {
	if n == 0 {
		return true
	}
	for i := 0; i < len(flowerbed); i++ {
		if flowerbed[i] == 0 {
			prev := (i == 0) || (flowerbed[i-1] == 0)
			next := (i == len(flowerbed)-1) || (flowerbed[i+1] == 0)

			if prev && next {
				flowerbed[i] = 1
				n--
				if n == 0 {
					return true
				}
			}
		}
	}
	return n == 0
}

func KidsWithCandies(candies []int, extraCandies int) []bool {
	max := 0
	for _, c := range candies {
		if c > max {
			max = c
		}
	}
	var result = make([]bool, len(candies))
	for i, c := range candies {
		if c+extraCandies >= max {
			result[i] = true
		}
	}
	return result
}

func MaxOperations(nums []int, k int) int {
	sort.Ints(nums)
	i := 0
	j := len(nums) - 1
	result := 0
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

func LongestOnes(nums []int, k int) int {
	left := 0
	zeros := 0
	maxLen := 0

	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			zeros++
		}

		for zeros > k {
			if nums[left] == 0 {
				zeros--
			}
			left++
		}

		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
}

func LongestSubarray(nums []int) int {
	left := 0
	zeros := 0
	maxLen := 0

	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			zeros++
		}

		for zeros > 1 {
			if nums[left] == 0 {
				zeros--
			}
			left++
		}

		maxLen = max(maxLen, right-left)
	}
	return maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
