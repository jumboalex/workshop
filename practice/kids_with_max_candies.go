package main

func kidsWithCandies(candies []int, extraCandies int) []bool {
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
