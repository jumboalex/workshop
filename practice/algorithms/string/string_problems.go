package stringproblems

func GcdOfStrings(str1 string, str2 string) string {
	prefix := str1
	if len(str1) > len(str2) {
		prefix = str2
	}

	for i := len(prefix); i >= 1; i-- {
		if len(str1)%i > 0 {
			prefix = prefix[:i-1]
			continue
		}
		if len(str2)%i > 0 {
			prefix = prefix[:i-1]
			continue
		}
		foundStr1 := true
		for j := 0; j < len(str1); j += i {
			if str1[j:j+i] != prefix {
				foundStr1 = false
				break
			}
		}
		foundStr2 := true
		for j := 0; j < len(str2); j += i {
			if str2[j:j+i] != prefix {
				foundStr2 = false
				break
			}
		}
		if foundStr1 && foundStr2 {
			return prefix
		} else {
			prefix = prefix[:i-1]
		}
	}
	return ""
}

func MergeAlternately(word1 string, word2 string) string {
	var result []byte
	for i := 0; i < len(word1) || i < len(word2); i++ {
		if i < len(word1) {
			result = append(result, word1[i])
		}
		if i < len(word2) {
			result = append(result, word2[i])
		}
	}
	return string(result)
}

func AddBinary(a string, b string) string {
	var result []byte

	carry := 0
	i := len(a) - 1
	j := len(b) - 1
	for i >= 0 && j >= 0 {
		digitA := int(a[i] - '0')
		digitB := int(b[j] - '0')
		sum := digitA + digitB + carry
		result = append([]byte{byte(sum%2 + '0')}, result...)
		carry = sum / 2
		i--
		j--
	}
	for i >= 0 {
		digitA := int(a[i] - '0')
		sum := digitA + carry
		result = append([]byte{byte(sum%2 + '0')}, result...)
		carry = sum / 2
		i--
	}
	for j >= 0 {
		digitB := int(b[j] - '0')
		sum := digitB + carry
		result = append([]byte{byte(sum%2 + '0')}, result...)
		carry = sum / 2
		j--
	}
	if carry > 0 {
		result = append([]byte{'1'}, result...)
	}
	return string(result)
}

func MultiplyString(num1 string, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}
	result := make([]int, len(num1)+len(num2))
	for i := len(num1) - 1; i >= 0; i-- {
		for j := len(num2) - 1; j >= 0; j-- {
			d := int(num1[i]-'0') * int(num2[j]-'0')
			s := d + result[i+j+1]

			result[i+j+1] = s % 10
			result[i+j] += s / 10
		}
	}
	if result[0] == 0 {
		result = result[1:]
	}
	bytes := make([]byte, len(result))
	for i, v := range result {
		bytes[i] = byte(v + '0')
	}
	return string(bytes)
}

func MaxVowels(s string, k int) int {
	vowels := map[byte]bool{
		'a': true, 'e': true, 'i': true, 'o': true, 'u': true,
	}

	currentVowels := 0
	for i := 0; i < k; i++ {
		if vowels[s[i]] {
			currentVowels++
		}
	}

	maxVowels := currentVowels
	for i := k; i < len(s); i++ {
		if vowels[s[i]] {
			currentVowels++
		}
		if vowels[s[i-k]] {
			currentVowels--
		}
		if currentVowels > maxVowels {
			maxVowels = currentVowels
		}
	}

	return maxVowels
}

func LengthOfLongestSubstringKDistinct(s string, k int) int {
	if k == 0 || len(s) == 0 {
		return 0
	}

	charCount := make(map[byte]int)
	left := 0
	maxLen := 0

	for right := 0; right < len(s); right++ {
		charCount[s[right]]++

		for len(charCount) > k {
			charCount[s[left]]--
			if charCount[s[left]] == 0 {
				delete(charCount, s[left])
			}
			left++
		}

		maxLen = max(maxLen, right-left+1)
	}

	return maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
