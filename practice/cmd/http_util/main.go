package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	fmt.Println("http utils")

	// Create http.Client with timeouts
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://jsonplaceholder.typicode.com/posts/1"

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add custom headers
	req.Header.Set("User-Agent", "Go-HTTP-Client/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("X-Custom-Header", "custom-value")

	// Execute request with client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Handle status codes
	if resp.StatusCode >= 400 {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(body))

	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("GET Post:", post.Body)

	// POST request example
	fmt.Println("\n--- POST Request Example ---")

	newPost := Post{
		UserID: 1,
		Title:  "New Post Title",
		Body:   "This is the body of the new post",
	}

	// Marshal post data to JSON
	postData, err := json.Marshal(newPost)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create POST request with body
	postURL := "https://jsonplaceholder.typicode.com/posts"
	postReq, err := http.NewRequest("POST", postURL, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set headers for POST
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Accept", "application/json")

	// Execute POST request
	postResp, err := client.Do(postReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer postResp.Body.Close()

	fmt.Println("POST Response Status:", postResp.Status)

	// Handle POST status codes
	switch {
	case postResp.StatusCode >= 200 && postResp.StatusCode < 300:
		fmt.Println("Success: POST request completed successfully")
	case postResp.StatusCode >= 400 && postResp.StatusCode < 500:
		fmt.Printf("Client Error: status code %d\n", postResp.StatusCode)
		return
	case postResp.StatusCode >= 500:
		fmt.Printf("Server Error: status code %d\n", postResp.StatusCode)
		return
	}

	// Read POST response
	postRespBytes, err := io.ReadAll(postResp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var createdPost Post
	err = json.Unmarshal(postRespBytes, &createdPost)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created Post ID:", createdPost.ID)
	fmt.Println("Created Post Title:", createdPost.Title)

}
