package main

func lengthOfLongestSubstringKDistinct(s string, k int) int {
	longestLength := 0
	start := 0
	end := 0
	charCount := make(map[byte]int)
	// sliding window
	for end < len(s) {
		charCount[s[end]]++

		for len(charCount) > k {
			charCount[s[start]]--
			if charCount[s[start]] == 0 {
				delete(charCount, s[start])
			}
			start++
		}
		if end-start+1 > longestLength {
			longestLength = end - start + 1
		}
		end++
	}
	return longestLength
}
