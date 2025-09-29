package arrayproblems

func RemoveDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}

	slow := 1
	fast := 1
	for ; fast < len(nums); fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
