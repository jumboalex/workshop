package main

// To execute Go code, please declare a func main() in a package "main"

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("Hello, World!")
	}
}

type WordFreq struct {
	token string
	count int
}

func concordance(filename string, n int) {
	dic := make(map[string]int)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		for _, w := range words {
			w := strings.ToLower(w)
			if _, ok := dic[w]; ok {
				dic[w]++
			} else {
				dic[w] = 1
			}
		}
	}

	wordFreq := []WordFreq{}
	for k, v := range dic {
		wordFreq = append(wordFreq, WordFreq{token: k, count: v})
	}

	slices.Sort(wordFreq, func(i, j int) bool {
		return wordFreq[i].count > wordFreq[j].count
	})

	for i := 0; i < n; i++ {
		fmt.Println(wordFreq[i].token, "   ", wordFreq[i].count)
	}
}

// how are you doing?
// Doing well!

// Your previous Plain Text content is preserved below:

// /* """
// Imagine you are given a plain text, ASCII, English document, something like a book from Project Gutenberg. Develop a concordance for the book -- the number of times each word appears -- and then print the top N most frequent words and how many times they occur. N can be either hardcoded or a parameter.

// === Example call and output ===
// $ ./concordance book.txt 10

// Output
// ------
// the     8230
// and     5067
// of      4139
// to      3651
// a       3017
// in      2659
// it      2082
// his     2008
// i       1972
// that    1950

// """ */

/*
	"""

A group of centaurs (mythical half-human, half-horse creatures) all sign
up for Facebook accounts at the same time. They immediately start
sending each other friend requests, in accordance with the ancient rules
that have governed centaur friendship since the dawn of time:

1) A centaur will only send a friend request to another centaur if the
recipient is at least (X/2 + 7) of the sender's age. For example, a 200-year
old centaur can only send friend requests to centaurs that are at least 107
years old.
2) A centaur will not send a friend request to another centaur that is older
than it is.
3) A centaur over 100 years old will not send a friend request to a recipient
under 100 years old. But centaurs under 100 years old can friend each other.
4) If any of the conditions for sending a friend request are not met, no
friend request will be sent.

Write a function that, given an array of centaur ages, returns an integer
of the total number of friend requests that the group of centaurs will send
to each other.

Examples:
count_all_friend_requests([120, 110])  => 1
# Friend requests          1    0
count_all_friend_requests([120, 110, 99]) => 1
# Friend requests          1    0    0
count_all_friend_requests([120, 45, 230, 400, 88, 300, 101]) => 4
# Friend requests          1    0   0    2    0   1    0
count_all_friend_requests([120, 45, 55, 230, 400, 88, 300, 101]) => 6
# Friend requests          1    0   1   0    2    1   1    0

"""
*/
func countRequest(centaurs []int) int {
	result := 0
	for i := 0; i < len(centaurs); i++ {
		for j := 0; j < len(centaurs); j++ {
			if i == j {
				continue
			}
			if centaurs[j] >= centaurs[i]/2+7 && centaurs[j] < centaurs[i] {
				result++
				continue
			}

		}
	}
	return result
}
