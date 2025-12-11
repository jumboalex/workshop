package main

import "fmt"

func main() {
	fmt.Println(findMaxAverage([]int{1, 12, -5, -6, 50, 3}, 4))
}

func findMaxAverage(nums []int, k int) float64 {
	start := 0
	end := k - 1

	maxAvg := -10001.0

	sum := 0
	for i := start; i <= end; i++ {
		sum += nums[i]
	}
	avg := float64(sum) / float64(k)
	if avg > maxAvg {
		maxAvg = avg
	}
	fmt.Println(sum, maxAvg)
	for end < len(nums) {
		end++

		if end == len(nums) {
			break
		}
		sum += nums[end]
		sum -= nums[start]
		start++
		avg := float64(sum) / float64(k)
		if avg > maxAvg {
			maxAvg = avg
		}
	}
	return maxAvg
}

func getAverages(nums []int, k int) []int {
	result := make([]int, len(nums))
	for i := range result {
		result[i] = -1
	}

	if 2*k+1 > len(nums) {
		return result
	}

	if k == 0 {
		return nums
	}

	// Build prefix sum array
	prefix := make([]int64, len(nums)+1)
	for i := 0; i < len(nums); i++ {
		prefix[i+1] = prefix[i] + int64(nums[i])
	}

	// Calculate k-radius averages using prefix sums
	for i := k; i < len(nums)-k; i++ {
		// Sum of window [i-k, i+k] = prefix[i+k+1] - prefix[i-k]
		windowSum := prefix[i+k+1] - prefix[i-k]
		result[i] = int(windowSum / int64(2*k+1))
	}
	return result
}
