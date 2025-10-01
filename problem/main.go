package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

// Your company runs a personal finance app that helps its users track how they spend their money. Your goal, for this problem, is to identify recurring subscriptions users have so that they may cancel unused ones.

// You have been provided a CSV file with one user's transactions. Each row corresponds to one transaction and contains the timestamp the transaction occurred, formatted as an ISO-8601 string. Find all recurring charges, then print the merchant, amount, and interval.

// Example output:
// "OrangeNews: $10.00 / week"

type Transaction struct {
	//created_at,merchant_name,amount,currency
	CreatedAt    time.Time
	MerchantName string
	Amount       int64
	Currency     string
}

func main() {
	fmt.Println("user subscription detect")

	filename := "testdata.csv"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// skip header
	_, err = reader.Read()
	if err != nil {
		fmt.Println(err)
	}

	var transactions []Transaction
	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Println(err)
		}

		t, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			fmt.Println(err)
		}
		amount, err := strconv.Atoi(record[2])
		if err != nil {
			fmt.Println(err)
		}

		transaction := Transaction{
			CreatedAt:    t,
			MerchantName: record[1],
			Amount:       int64(amount),
			Currency:     record[3],
		}

		transactions = append(transactions, transaction)
	}

	merchantReport := make(map[string]map[time.Time]int64)
	for _, transact := range transactions {
		if _, ok := merchantReport[transact.MerchantName]; !ok {
			subsInterval := make(map[time.Time]int64)
			subsInterval[transact.CreatedAt] = transact.Amount
			merchantReport[transact.MerchantName] = subsInterval
		} else {
			merchantReport[transact.MerchantName][transact.CreatedAt] = transact.Amount
		}
	}

	for m, v := range merchantReport {
		times := []time.Time{}
		for k := range v {
			times = append(times, k)
		}
		sort.Slice(times, func(i, j int) bool {
			return times[i].Before(times[j])
		})
		left := 0
		right := 0
		for right < len(times) {
			d := times[right].Sub(times[left])
			for d.Hours() < 7*24 {
				// adding amount
			}

		}
	}

}
