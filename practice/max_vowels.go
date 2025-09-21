package main

func maxVowels(s string, k int) int {
	vowels := map[byte]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}
	/*
		maxCount := 0
		for i := 0; i < len(s)-k; i++ {
			substr := s[i : i+k]
			fmt.Println("substr:", substr)
			count := 0
			for j := 0; j < k; j++ {
				fmt.Println("i:", i, " j:", j, " substr[j]:", substr[j])
				if _, ok := vowels[substr[j]]; ok {
					count++
				}
			}
			if count > maxCount {
				maxCount = count
			}
		}
	*/

	// sliding window
	count := 0
	for i := 0; i < k; i++ {
		if vowels[s[i]] {
			count++
		}
	}
	maxCount := count
	for i := k; i < len(s); i++ {
		if vowels[s[i]] {
			count++
		}
		if vowels[s[i-k]] {
			count--
		}
		if count > maxCount {
			maxCount = count
		}
	}

	return maxCount
}
