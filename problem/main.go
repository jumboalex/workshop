package main

import (
	"encoding/csv"
	"fmt"
	"io"
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
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	var transactions []Transaction
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		t, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			continue
		}
		amount, err := strconv.Atoi(record[2])
		if err != nil {
			continue
		}

		transaction := Transaction{
			CreatedAt:    t,
			MerchantName: record[1],
			Amount:       int64(amount),
			Currency:     record[3],
		}

		transactions = append(transactions, transaction)
	}

	merchantReport := make(map[string][]Transaction)
	for _, transact := range transactions {
		merchantReport[transact.MerchantName] = append(merchantReport[transact.MerchantName], transact)
	}

	for merchant, merchantTransactions := range merchantReport {
		if len(merchantTransactions) < 2 {
			continue
		}

		sort.Slice(merchantTransactions, func(i, j int) bool {
			return merchantTransactions[i].CreatedAt.Before(merchantTransactions[j].CreatedAt)
		})

		// Calculate intervals between consecutive transactions
		intervals := make(map[int]int) // interval in days -> count
		for i := 1; i < len(merchantTransactions); i++ {
			hours := merchantTransactions[i].CreatedAt.Sub(merchantTransactions[i-1].CreatedAt).Hours()
			days := int(hours/24 + 0.5) // Round to nearest day
			intervals[days]++
		}

		// Find the most common interval (if it appears more than once)
		var mostCommonInterval int
		var maxCount int
		for interval, count := range intervals {
			if count > maxCount {
				maxCount = count
				mostCommonInterval = interval
			}
		}

		// Only consider it recurring if the interval is weekly (7 days)
		if mostCommonInterval == 7 && maxCount >= 1 {
			amount := merchantTransactions[0].Amount
			amountStr := fmt.Sprintf("$%.2f", float64(amount)/100)
			fmt.Printf("%s: %s / week\n", merchant, amountStr)
		}
	}

}
