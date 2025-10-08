package stringproblems

import "testing"

func TestGcdOfStrings(t *testing.T) {
	tests := []struct {
		name string
		str1 string
		str2 string
		want string
	}{
		{"example 1", "ABCABC", "ABC", "ABC"},
		{"example 2", "ABABAB", "ABAB", "AB"},
		{"no gcd", "LEET", "CODE", ""},
		{"same strings", "ABC", "ABC", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GcdOfStrings(tt.str1, tt.str2)
			if got != tt.want {
				t.Errorf("GcdOfStrings(%q, %q) = %q, want %q", tt.str1, tt.str2, got, tt.want)
			}
		})
	}
}

func TestMergeAlternately(t *testing.T) {
	tests := []struct {
		name  string
		word1 string
		word2 string
		want  string
	}{
		{"equal length", "abc", "pqr", "apbqcr"},
		{"first longer", "ab", "pqrs", "apbqrs"},
		{"second longer", "abcd", "pq", "apbqcd"},
		{"empty first", "", "abc", "abc"},
		{"empty second", "abc", "", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeAlternately(tt.word1, tt.word2)
			if got != tt.want {
				t.Errorf("MergeAlternately(%q, %q) = %q, want %q", tt.word1, tt.word2, got, tt.want)
			}
		})
	}
}

func TestAddBinary(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"simple", "11", "1", "100"},
		{"example 1", "1010", "1011", "10101"},
		{"zeros", "0", "0", "0"},
		{"carry", "1111", "1111", "11110"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddBinary(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("AddBinary(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiplyString(t *testing.T) {
	tests := []struct {
		name string
		num1 string
		num2 string
		want string
	}{
		{"simple", "2", "3", "6"},
		{"with zero", "123", "0", "0"},
		{"large numbers", "123", "456", "56088"},
		{"single digits", "9", "9", "81"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultiplyString(tt.num1, tt.num2)
			if got != tt.want {
				t.Errorf("MultiplyString(%q, %q) = %q, want %q", tt.num1, tt.num2, got, tt.want)
			}
		})
	}
}

func TestMaxVowels(t *testing.T) {
	tests := []struct {
		name string
		s    string
		k    int
		want int
	}{
		{"example 1", "abciiidef", 3, 3},
		{"example 2", "aeiou", 2, 2},
		{"no vowels", "rhythms", 4, 0},
		{"all vowels", "aeiou", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxVowels(tt.s, tt.k)
			if got != tt.want {
				t.Errorf("MaxVowels(%q, %d) = %v, want %v", tt.s, tt.k, got, tt.want)
			}
		})
	}
}

func TestLengthOfLongestSubstringKDistinct(t *testing.T) {
	tests := []struct {
		name string
		s    string
		k    int
		want int
	}{
		{"example 1", "eceba", 2, 3},
		{"example 2", "aa", 1, 2},
		{"k is 0", "abc", 0, 0},
		{"empty string", "", 2, 0},
		{"longer string", "abcadcacacaca", 3, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LengthOfLongestSubstringKDistinct(tt.s, tt.k)
			if got != tt.want {
				t.Errorf("LengthOfLongestSubstringKDistinct(%q, %d) = %v, want %v", tt.s, tt.k, got, tt.want)
			}
		})
	}
}

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name string
		s    string
		p    string
		want []int
	}{
		{"example 1", "cbaebabacd", "abc", []int{0, 6}},
		{"example 2", "abab", "ab", []int{0, 1, 2}},
		{"no anagrams", "hello", "world", []int{}},
		{"p longer than s", "ab", "abc", []int{}},
		{"single char", "aaa", "a", []int{0, 1, 2}},
		{"empty s", "", "abc", []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindAnagrams(tt.s, tt.p)
			if len(got) != len(tt.want) {
				t.Errorf("FindAnagrams(%q, %q) = %v, want %v", tt.s, tt.p, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("FindAnagrams(%q, %q) = %v, want %v", tt.s, tt.p, got, tt.want)
					return
				}
			}
		})
	}
}
