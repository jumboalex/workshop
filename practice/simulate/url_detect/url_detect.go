package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FetchedResponse struct {
	Urls    []string `json:"urls"`
	Message string   `json:"message"`
}

func main() {
	fmt.Println("detect the exit of url fetching")

	// bfs
	visited := make(map[string]bool)
	q := []string{"http://example.com"}
	for len(q) > 0 {
		url := q[0]
		q = q[1:]

		if visited[url] {
			continue
		}
		visited[url] = true
		urls, shouldReturn := fetchURLResponse(url)
		if shouldReturn {
			return
		}

		for _, u := range urls {
			if !visited[u] {
				q = append(q, u)
			}
		}
	}
}

const retryAttempts = 3

func fetchURLResponse(url string) ([]string, bool) {
	client := &http.Client{Timeout: 15 * time.Second}
	for i := 0; i < retryAttempts; i++ {
		response, err := client.Get(url)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return []string{}, false
		}

		if response.StatusCode >= 500 {
			// need to retry
			fmt.Println("Server error, need to retry:", response.Status)
			time.Sleep(500 * time.Millisecond)
			response.Body.Close()
			continue
		}
		if response.StatusCode >= 400 {
			// client error, do not retry
			fmt.Println("Client error, do not retry:", response.Status)
			response.Body.Close()
			return []string{}, false
		}
		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			response.Body.Close()
			return []string{}, false
		}
		response.Body.Close()
		var fetchedResponse FetchedResponse
		err = json.Unmarshal(bytes, &fetchedResponse)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return []string{}, false
		}

		if fetchedResponse.Message == "exit" {
			fmt.Println("Exit message received. Stopping further requests.")
			return []string{}, true
		} else {
			return fetchedResponse.Urls, false
		}

	}

	return []string{}, false
}
