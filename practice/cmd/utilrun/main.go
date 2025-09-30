package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type CarData struct {
	Model        string  `json:"model"`
	MPG          float64 `json:"mpg"`
	Cylinders    int     `json:"cyl"`
	Displacement float64 `json:"disp"`
	Drat         float64 `json:"drat"`
	Weight       float64 `json:"wt"`
	Qsec         float64 `json:"qsec"`
	VS           int     `json:"vs"`
	AM           int     `json:"am"`
	Gear         int     `json:"gear"`
	Carb         int     `json:"carb"`
}

func main() {
	fmt.Println("test utils")

	readJSONByLine()

	fmt.Println("================")

	readJSONWholeFile()

	fmt.Println("=================")

	filepath := "mtcars-parquet-1.json"
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}

	var carMap []map[string]any
	err = json.Unmarshal(bytes, &carMap)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range carMap {
		for key, value := range v {
			fmt.Printf("%s: %v, ", key, value)
		}
		fmt.Println()
	}
}

func readJSONWholeFile() {
	filepath := "mtcars-parquet-1.json"
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}

	var cars []CarData
	err = json.Unmarshal(bytes, &cars)
	if err != nil {
		fmt.Println(err)
	}
	for _, car := range cars {
		fmt.Println(car)
	}
}

func readJSONByLine() {
	filepath := "mtcars-parquet.json"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var car CarData
		err := json.Unmarshal(scanner.Bytes(), &car)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		fmt.Println(car)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
