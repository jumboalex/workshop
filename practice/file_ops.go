package main

import (
	"bufio"
	"encoding/csv"
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

	content, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(content))

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
