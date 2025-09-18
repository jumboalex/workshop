package main

import (
	"bufio"
	"fmt"
	"os"
)

func readFile() {
	file, err := os.Open("customers-100.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

}
