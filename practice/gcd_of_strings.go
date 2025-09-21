package main

import "fmt"

func gcdOfStrings(str1 string, str2 string) string {
	prefix := str1
	if len(str1) > len(str2) {
		prefix = str2
	}

	fmt.Println("prefix:", prefix)
	for i := len(prefix); i >= 1; i-- {
		fmt.Println("i:", i)
		if len(str1)%i > 0 {
			prefix = prefix[:i-1]
			continue
		}
		if len(str2)%i > 0 {
			prefix = prefix[:i-1]
			continue
		}
		fmt.Println("prefix:", prefix)
		foundStr1 := true
		for j := 0; j < len(str1); j += i {
			fmt.Println("j:", j, " str1[j:j+i]:", str1[j:j+i], " prefix:", prefix)
			if str1[j:j+i] != prefix {
				foundStr1 = false
				break
			}
		}
		foundStr2 := true
		for j := 0; j < len(str2); j += i {
			fmt.Println("j:", j, " str2[j:j+i]:", str2[j:j+i], " prefix:", prefix)
			if str2[j:j+i] != prefix {
				foundStr2 = false
				break
			}
		}
		if foundStr1 && foundStr2 {
			return prefix
		} else {
			prefix = prefix[:i-1]
			fmt.Println("new prefix:", prefix)
		}
	}
	return ""
}
