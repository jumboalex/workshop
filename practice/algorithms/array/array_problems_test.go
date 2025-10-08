package arrayproblems

import (
	"reflect"
	"testing"
)

func TestPlusOne(t *testing.T) {
	tests := []struct {
		name   string
		digits []int
		want   []int
	}{
		{"single digit", []int{9}, []int{1, 0}},
		{"no carry", []int{1, 2, 3}, []int{1, 2, 4}},
		{"carry", []int{1, 9, 9}, []int{2, 0, 0}},
		{"all nines", []int{9, 9, 9}, []int{1, 0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]int, len(tt.digits))
			copy(input, tt.digits)
			got := PlusOne(input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlusOne(%v) = %v, want %v", tt.digits, got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		wantLen  int
		wantNums []int
	}{
		{"empty", []int{}, 0, []int{}},
		{"single", []int{1}, 1, []int{1}},
		{"no duplicates", []int{1, 2, 3}, 3, []int{1, 2, 3}},
		{"with duplicates", []int{1, 1, 2}, 2, []int{1, 2}},
		{"multiple duplicates", []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}, 5, []int{0, 1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)
			got := RemoveDuplicates(nums)
			if got != tt.wantLen {
				t.Errorf("RemoveDuplicates() length = %v, want %v", got, tt.wantLen)
			}
			if !reflect.DeepEqual(nums[:got], tt.wantNums) {
				t.Errorf("RemoveDuplicates() nums = %v, want %v", nums[:got], tt.wantNums)
			}
		})
	}
}

func TestCanPlaceFlowers(t *testing.T) {
	tests := []struct {
		name      string
		flowerbed []int
		n         int
		want      bool
	}{
		{"can place", []int{1, 0, 0, 0, 1}, 1, true},
		{"cannot place", []int{1, 0, 0, 0, 1}, 2, false},
		{"empty bed", []int{0, 0, 0}, 2, true},
		{"zero flowers", []int{1, 0, 1}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flowerbed := make([]int, len(tt.flowerbed))
			copy(flowerbed, tt.flowerbed)
			got := CanPlaceFlowers(flowerbed, tt.n)
			if got != tt.want {
				t.Errorf("CanPlaceFlowers(%v, %d) = %v, want %v", tt.flowerbed, tt.n, got, tt.want)
			}
		})
	}
}

func TestKidsWithCandies(t *testing.T) {
	tests := []struct {
		name         string
		candies      []int
		extraCandies int
		want         []bool
	}{
		{"example 1", []int{2, 3, 5, 1, 3}, 3, []bool{true, true, true, false, true}},
		{"example 2", []int{4, 2, 1, 1, 2}, 1, []bool{true, false, false, false, false}},
		{"all equal", []int{1, 1, 1}, 0, []bool{true, true, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KidsWithCandies(tt.candies, tt.extraCandies)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KidsWithCandies(%v, %d) = %v, want %v", tt.candies, tt.extraCandies, got, tt.want)
			}
		})
	}
}

func TestMaxOperations(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"example 1", []int{1, 2, 3, 4}, 5, 2},
		{"example 2", []int{3, 1, 3, 4, 3}, 6, 1},
		{"no pairs", []int{1, 2, 3}, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)
			got := MaxOperations(nums, tt.k)
			if got != tt.want {
				t.Errorf("MaxOperations(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

func TestLongestOnes(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"example 1", []int{1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0}, 2, 6},
		{"example 2", []int{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1}, 3, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestOnes(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("LongestOnes(%v, %d) = %v, want %v", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

func TestLongestSubarray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"example 1", []int{1, 1, 0, 1}, 3},
		{"example 2", []int{0, 1, 1, 1, 0, 1, 1, 0, 1}, 5},
		{"all ones", []int{1, 1, 1}, 2},
		{"all zeros", []int{0, 0, 0}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestSubarray(tt.nums)
			if got != tt.want {
				t.Errorf("LongestSubarray(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestFindDiagonalOrder(t *testing.T) {
	tests := []struct {
		name string
		mat  [][]int
		want []int
	}{
		{
			"3x3 matrix",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 4, 2, 7, 5, 3, 8, 6, 9},
		},
		{
			"2x2 matrix",
			[][]int{{1, 2}, {3, 4}},
			[]int{1, 3, 2, 4},
		},
		{
			"1x3 matrix",
			[][]int{{1, 2, 3}},
			[]int{1, 2, 3},
		},
		{
			"3x1 matrix",
			[][]int{{1}, {2}, {3}},
			[]int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findDiagonalOrder(tt.mat)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDiagonalOrder(%v) = %v, want %v", tt.mat, got, tt.want)
			}
		})
	}
}

func TestParseIPv6(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []int
	}{
		{
			"full address",
			"1234:abcd:0000:0000:0102:0000:0000:fffe",
			[]int{18, 52, 171, 205, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 255, 254},
		},
		{
			"compressed middle",
			"2001:db8::1",
			[]int{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			"localhost",
			"::1",
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			"compressed end",
			"1234::",
			[]int{18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"all zeros",
			"::",
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"IPv4 embedded",
			"::ffff:192.168.1.1",
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 1, 1},
		},
		{
			"compressed complex",
			"2001:db8:85a3::8a2e:370:7334",
			[]int{32, 1, 13, 184, 133, 163, 0, 0, 0, 0, 138, 46, 3, 112, 115, 52},
		},
		{
			"IPv4 mapping common",
			"0:0:0:0:0:0:13.1.68.3",
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 1, 68, 3},
		},
		{
			"empty string",
			"",
			nil,
		},
		{
			"invalid hex",
			"gggg::1",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseIPv6(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIPv6(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
