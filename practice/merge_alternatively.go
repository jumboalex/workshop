package main

func mergeAlternately(word1 string, word2 string) string {
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
