package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// --- Configuration ---
const (
	// Replace this with the actual URL when provided.
	SatAPIURL = "http://satellite-api.example.com"
	// A reasonable timeout for any single request over a slow, unreliable link.
	RequestTimeout = 30 * time.Second
	// Max retries for transient errors (connection, 5xx server errors).
	MaxRetries = 5
)

// --- Data Models ---

// SensorStatus is an enum-like type for the sensor status.
type SensorStatus string

const (
	StatusInitializing SensorStatus = "INITIALIZING"
	StatusActive       SensorStatus = "ACTIVE"
	StatusFailed       SensorStatus = "FAILED"
	StatusRestarting   SensorStatus = "RESTARTING"
	StatusTerminating  SensorStatus = "TERMINATING"
)

// Sensor represents the Sensor object data model.
type Sensor struct {
	ID          int          `json:"id"`
	Frequency   int          `json:"frequency"`
	Status      SensorStatus `json:"status"`
	Measurement *float64     `json:"measurement"` // Use pointer for nullability
}

// SensorCreateRequest is the payload for POST /sensors.
type SensorCreateRequest struct {
	Frequency int `json:"frequency"`
}

// --- Custom Error ---

// SatelliteAPIError represents an error from the API, usually non-retriable 4xx or a
// final 5xx error after all retries.
type SatelliteAPIError struct {
	StatusCode int
	Message    string
}

func (e *SatelliteAPIError) Error() string {
	return fmt.Sprintf("API Error (Status %d): %s", e.StatusCode, e.Message)
}

// --- Interface Implementation ---

// SatelliteInterface manages communication with the satellite API.
type SatelliteInterface struct {
	BaseURL string
	Client  *retryablehttp.Client
}

// NewSatelliteInterface creates a new interface client configured for resilience.
func NewSatelliteInterface(baseURL string) *SatelliteInterface {
	// Configure the retryable HTTP client
	client := retryablehttp.NewClient()
	client.HTTPClient.Timeout = RequestTimeout // Timeout for each request attempt
	client.RetryMax = MaxRetries               // Max number of retries
	// Use an exponential backoff strategy (default behavior)
	client.Backoff = retryablehttp.DefaultBackoff

	// Custom check for retriable status codes (default is 429 and 5xx)
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			// Do NOT retry on most 4xx client errors
			return false, nil
		}
		// Use the default retry logic for transient errors (connection errors, 5xx, 429)
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	return &SatelliteInterface{
		BaseURL: baseURL,
		Client:  client,
	}
}

// internalRequest executes a resilient HTTP request and handles API errors.
func (s *SatelliteInterface) internalRequest(method, endpoint string, reqBody interface{}, respData interface{}) error {
	url := s.BaseURL + endpoint

	// 1. Prepare Request Body (if any)
	var bodyReader io.Reader
	if reqBody != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		bodyReader = buf
	}

	// Create the request
	req, err := retryablehttp.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 2. Execute Request
	resp, err := s.Client.Do(req)
	if err != nil {
		// This includes errors after all retries have failed (timeout/connection issues)
		return fmt.Errorf("request failed after %d retries: %w", s.Client.RetryMax, err)
	}
	defer resp.Body.Close()

	// 3. Handle Status Codes
	if resp.StatusCode >= 400 {
		return &SatelliteAPIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Request to %s failed with status %d", endpoint, resp.StatusCode),
		}
	}

	// 4. Decode Response Data
	if respData != nil && resp.ContentLength != 0 {
		if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}

// GetSensorIDs GET /sensor-ids
func (s *SatelliteInterface) GetSensorIDs() ([]int, error) {
	fmt.Println("Attempting to GET /sensor-ids...")
	var ids []int
	if err := s.internalRequest(http.MethodGet, "/sensor-ids", nil, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// CreateSensor POST /sensors
func (s *SatelliteInterface) CreateSensor(frequency int) (*Sensor, error) {
	fmt.Printf("Attempting to POST /sensors with frequency: %d...\n", frequency)
	req := SensorCreateRequest{Frequency: frequency}
	sensor := &Sensor{}
	if err := s.internalRequest(http.MethodPost, "/sensors", req, sensor); err != nil {
		return nil, err
	}
	return sensor, nil
}

// GetSensor GET /sensors/<id>
// Automatically retries if sensor status is INITIALIZING or RESTARTING (up to 30 seconds)
func (s *SatelliteInterface) GetSensor(id int) (*Sensor, error) {
	endpoint := fmt.Sprintf("/sensors/%d", id)
	fmt.Printf("Attempting to GET %s...\n", endpoint)

	maxWait := 30 * time.Second
	pollInterval := 2 * time.Second
	deadline := time.Now().Add(maxWait)

	for {
		sensor := &Sensor{}
		if err := s.internalRequest(http.MethodGet, endpoint, nil, sensor); err != nil {
			return nil, err
		}

		// Return immediately if active
		if sensor.Status == StatusActive {
			return sensor, nil
		}

		// Return error for terminal states
		if sensor.Status == StatusFailed || sensor.Status == StatusTerminating {
			return nil, fmt.Errorf("sensor %d is in terminal state: %s", id, sensor.Status)
		}

		// Check timeout
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for sensor %d to become ACTIVE (current status: %s)", id, sensor.Status)
		}

		// Retry for INITIALIZING or RESTARTING states
		fmt.Printf("   Sensor status: %s, retrying...\n", sensor.Status)
		time.Sleep(pollInterval)
	}
}

// --- Helper Functions ---

func handleError(op string, err error) {
	if apiErr, ok := err.(*SatelliteAPIError); ok {
		fmt.Printf("\n❌ A non-retriable API error occurred during %s:\n", op)
		fmt.Printf("   Status: %d\n", apiErr.StatusCode)
		fmt.Printf("   Message: %s\n", apiErr.Message)
		// Optionally print the error body
		// fmt.Printf("   Body: %+v\n", apiErr.Body)
	} else {
		fmt.Printf("\n❌ A critical error occurred during %s after all retries failed:\n", op)
		fmt.Printf("   Error: %v\n", err)
	}
}
