package main

import (
	"strconv"
	"strings"
)

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
	resultStr := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		resultStr[i] = strconv.Itoa(result[i])
	}
	if resultStr[0] == "0" {
		resultStr = resultStr[1:]
	}
	return strings.Join(resultStr, "")
}
