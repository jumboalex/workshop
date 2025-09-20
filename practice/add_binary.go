package main

import "fmt"

func addBinary(a string, b string) string {
	fmt.Println("test addBinary")
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
