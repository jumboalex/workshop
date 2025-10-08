package arrayproblems

import (
	"fmt"
	"sort"
)

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

func findDiagonalOrder(mat [][]int) []int {
	result := []int{}
	m := len(mat)
	n := len(mat[0])

	for d := 0; d < m+n-1; d++ {
		var startRow int
		var startCol int

		if d < n {
			startRow = 0
			startCol = d
		} else {
			startRow = d - n + 1
			startCol = n - 1
		}

		for i, j := startRow, startCol; i < m && j >= 0; i, j = i+1, j-1 {
			result = append(result, mat[i][j])
		}
	}
	return result
}

func PrintDiagonalOrder(mat [][]int) {
	m := len(mat)
	n := len(mat[0])

	for d := 0; d < m+n-1; d++ {
		var startRow int
		var startCol int

		if d < n {
			startRow = 0
			startCol = d
		} else {
			startRow = d - n + 1
			startCol = n - 1
		}

		for i, j := startRow, startCol; i < m && j >= 0; i, j = i+1, j-1 {
			fmt.Print(mat[i][j], " ")
		}
	}
	fmt.Println()
}

func IPv4ToIPv6(ipv4 string) string {
	// Edge cases
	if ipv4 == "" {
		return ""
	}

	// Split by '.'
	parts := []string{}
	current := ""
	for i := 0; i < len(ipv4); i++ {
		if ipv4[i] == '.' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(ipv4[i])
		}
	}
	parts = append(parts, current)

	// Validate: must have exactly 4 parts
	if len(parts) != 4 {
		return ""
	}

	// Validate and convert each octet
	octets := make([]int, 4)
	for i, part := range parts {
		// Empty part
		if part == "" {
			return ""
		}

		// Leading zeros (except "0" itself)
		if len(part) > 1 && part[0] == '0' {
			return ""
		}

		// Convert to int
		num := 0
		for j := 0; j < len(part); j++ {
			if part[j] < '0' || part[j] > '9' {
				return ""
			}
			num = num*10 + int(part[j]-'0')
		}

		// Must be 0-255
		if num > 255 {
			return ""
		}

		octets[i] = num
	}

	// Convert to IPv6 format: ::ffff:a.b.c.d
	return fmt.Sprintf("::ffff:%d.%d.%d.%d", octets[0], octets[1], octets[2], octets[3])
}

func ParseIPv6(ipv6 string) []int {
	// Edge case: empty string
	if ipv6 == "" {
		return nil
	}

	// Check for embedded IPv4 (e.g., ::ffff:192.168.1.1)
	hasIPv4 := false
	ipv4Start := -1
	for i := 0; i < len(ipv6); i++ {
		if ipv6[i] == '.' {
			hasIPv4 = true
			// Find start of IPv4 part
			for j := i; j >= 0; j-- {
				if j == 0 || ipv6[j-1] == ':' {
					ipv4Start = j
					break
				}
			}
			break
		}
	}

	var ipv4Bytes []int
	hexPart := ipv6
	if hasIPv4 && ipv4Start >= 0 {
		ipv4Str := ipv6[ipv4Start:]
		hexPart = ipv6[:ipv4Start]
		if len(hexPart) > 0 && hexPart[len(hexPart)-1] == ':' {
			hexPart = hexPart[:len(hexPart)-1]
		}

		// Parse IPv4
		parts := splitString(ipv4Str, '.')
		if len(parts) != 4 {
			return nil
		}
		for _, part := range parts {
			num := parseDecimal(part)
			if num < 0 || num > 255 {
				return nil
			}
			ipv4Bytes = append(ipv4Bytes, num)
		}
	}

	// Handle :: compression
	hasDoubleColon := false
	doubleColonIdx := -1
	for i := 0; i < len(hexPart)-1; i++ {
		if hexPart[i] == ':' && hexPart[i+1] == ':' {
			if hasDoubleColon {
				return nil // Multiple :: not allowed
			}
			hasDoubleColon = true
			doubleColonIdx = i
		}
	}

	var groups []string
	var leftGroups, rightGroups []string

	if hasDoubleColon {
		left := hexPart[:doubleColonIdx]
		right := ""
		if doubleColonIdx+2 < len(hexPart) {
			right = hexPart[doubleColonIdx+2:]
		}

		if left != "" {
			leftGroups = splitString(left, ':')
		}
		if right != "" {
			rightGroups = splitString(right, ':')
		}
	} else {
		groups = splitString(hexPart, ':')
	}

	// Convert hex groups to bytes
	result := []int{}

	if hasDoubleColon {
		// Add left groups
		for _, g := range leftGroups {
			val := parseHex(g)
			if val < 0 {
				return nil
			}
			result = append(result, val>>8, val&0xff)
		}

		// Calculate zeros needed
		totalGroups := 8
		if hasIPv4 {
			totalGroups = 6 // IPv4 takes last 2 groups
		}
		zerosNeeded := totalGroups - len(leftGroups) - len(rightGroups)
		for i := 0; i < zerosNeeded*2; i++ {
			result = append(result, 0)
		}

		// Add right groups
		for _, g := range rightGroups {
			val := parseHex(g)
			if val < 0 {
				return nil
			}
			result = append(result, val>>8, val&0xff)
		}
	} else {
		expectedGroups := 8
		if hasIPv4 {
			expectedGroups = 6
		}
		if len(groups) != expectedGroups {
			return nil
		}

		for _, g := range groups {
			val := parseHex(g)
			if val < 0 {
				return nil
			}
			result = append(result, val>>8, val&0xff)
		}
	}

	// Add IPv4 bytes if present
	if hasIPv4 {
		result = append(result, ipv4Bytes...)
	}

	if len(result) != 16 {
		return nil
	}

	return result
}

func splitString(s string, sep byte) []string {
	if s == "" {
		return []string{}
	}
	parts := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	parts = append(parts, current)
	return parts
}

func parseHex(s string) int {
	if s == "" || len(s) > 4 {
		return -1
	}
	val := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		var digit int
		if c >= '0' && c <= '9' {
			digit = int(c - '0')
		} else if c >= 'a' && c <= 'f' {
			digit = int(c-'a') + 10
		} else if c >= 'A' && c <= 'F' {
			digit = int(c-'A') + 10
		} else {
			return -1
		}
		val = val*16 + digit
	}
	return val
}

func parseDecimal(s string) int {
	if s == "" {
		return -1
	}
	val := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return -1
		}
		val = val*10 + int(s[i]-'0')
	}
	return val
}
