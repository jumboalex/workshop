package arrayproblems

func PlusOne(digits []int) []int {
	borrow := 1
	for i := len(digits) - 1; i >= 0; i-- {
		digits[i] += borrow
		if digits[i] == 10 {
			digits[i] = 0
			borrow = 1
		} else {
			borrow = 0
		}
	}
	if borrow == 1 {
		digits = append([]int{1}, digits...)
	}
	return digits
}
