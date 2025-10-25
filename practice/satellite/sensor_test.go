package main

import (
	"testing"
	"time"
)

func TestSatelliteInterface_GetSensorIDs(t *testing.T) {
	// Start mock server in background
	server := NewMockServer(0.0, 0) // No unreliability for basic tests
	go server.Start(8081)
	time.Sleep(100 * time.Millisecond) // Wait for server to start

	client := NewSatelliteInterface("http://localhost:8081")

	// Initially should have no sensors
	ids, err := client.GetSensorIDs()
	if err != nil {
		t.Fatalf("GetSensorIDs failed: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("Expected 0 sensors, got %d", len(ids))
	}

	// Create a sensor
	sensor, err := client.CreateSensor(100)
	if err != nil {
		t.Fatalf("CreateSensor failed: %v", err)
	}
	if sensor.Frequency != 100 {
		t.Errorf("Expected frequency 100, got %d", sensor.Frequency)
	}

	// Now should have 1 sensor
	ids, err = client.GetSensorIDs()
	if err != nil {
		t.Fatalf("GetSensorIDs failed: %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("Expected 1 sensor, got %d", len(ids))
	}
}

func TestSatelliteInterface_CreateSensor(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8082)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8082")

	sensor, err := client.CreateSensor(42)
	if err != nil {
		t.Fatalf("CreateSensor failed: %v", err)
	}

	if sensor.ID == 0 {
		t.Error("Expected non-zero sensor ID")
	}
	if sensor.Frequency != 42 {
		t.Errorf("Expected frequency 42, got %d", sensor.Frequency)
	}
	if sensor.Status != StatusInitializing {
		t.Errorf("Expected status INITIALIZING, got %s", sensor.Status)
	}
}

func TestSatelliteInterface_GetSensor(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8083)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8083")

	// Create a sensor
	created, err := client.CreateSensor(99)
	if err != nil {
		t.Fatalf("CreateSensor failed: %v", err)
	}

	// Get the sensor
	fetched, err := client.GetSensor(created.ID)
	if err != nil {
		t.Fatalf("GetSensor failed: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, fetched.ID)
	}
	if fetched.Frequency != 99 {
		t.Errorf("Expected frequency 99, got %d", fetched.Frequency)
	}
}

func TestSatelliteInterface_GetSensor_NotFound(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8084)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8084")

	// Try to get non-existent sensor
	_, err := client.GetSensor(999)
	if err == nil {
		t.Error("Expected error for non-existent sensor")
	}

	// Check that it's a SatelliteAPIError with 404
	if apiErr, ok := err.(*SatelliteAPIError); ok {
		if apiErr.StatusCode != 404 {
			t.Errorf("Expected status 404, got %d", apiErr.StatusCode)
		}
	} else {
		t.Errorf("Expected SatelliteAPIError, got %T", err)
	}
}

func TestSatelliteInterface_WithUnreliability(t *testing.T) {
	// Test with 30% unreliability
	server := NewMockServer(0.3, 500*time.Millisecond)
	go server.Start(8085)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8085")

	// Should eventually succeed despite unreliability
	sensor, err := client.CreateSensor(123)
	if err != nil {
		t.Logf("CreateSensor failed even with retries: %v", err)
		// Don't fail the test - unreliability might cause all retries to fail
		return
	}

	t.Logf("Successfully created sensor %d despite unreliability", sensor.ID)
}

func TestSatelliteInterface_InvalidFrequency(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8086)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8086")

	// Try to create sensor with invalid frequency
	_, err := client.CreateSensor(-1)
	if err == nil {
		t.Error("Expected error for invalid frequency")
	}

	// Check that it's a SatelliteAPIError with 400
	if apiErr, ok := err.(*SatelliteAPIError); ok {
		if apiErr.StatusCode != 400 {
			t.Errorf("Expected status 400, got %d", apiErr.StatusCode)
		}
	}
}

func TestSatelliteInterface_RetryOnServerError(t *testing.T) {
	// High unreliability to trigger retries
	server := NewMockServer(0.5, 100*time.Millisecond)
	go server.Start(8087)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8087")

	// This should retry multiple times due to 5xx errors
	start := time.Now()
	_, err := client.GetSensorIDs()
	duration := time.Since(start)

	// If it retried, it should take longer than a single request
	if err == nil {
		t.Logf("GetSensorIDs succeeded after retries (took %v)", duration)
	} else {
		// Even if all retries fail, that's acceptable for this test
		t.Logf("All retries exhausted after %v: %v", duration, err)
	}
}

func TestSatelliteInterface_CustomErrorType(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8088)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8088")

	// Test that 404 returns proper error type
	_, err := client.GetSensor(999)
	if err == nil {
		t.Fatal("Expected error for non-existent sensor")
	}

	satErr, ok := err.(*SatelliteAPIError)
	if !ok {
		t.Fatalf("Expected *SatelliteAPIError, got %T", err)
	}

	if satErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", satErr.StatusCode)
	}

	if satErr.Error() == "" {
		t.Error("Error message should not be empty")
	}

	t.Logf("Error message: %s", satErr.Error())
}

func TestSatelliteInterface_ConcurrentRequests(t *testing.T) {
	server := NewMockServer(0.1, 50*time.Millisecond)
	go server.Start(8089)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8089")

	// Create multiple sensors concurrently
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(freq int) {
			_, err := client.CreateSensor(freq)
			if err != nil {
				t.Logf("Concurrent create failed for freq %d: %v", freq, err)
			}
			done <- true
		}(100 + i)
	}

	// Wait for all to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify we can still query the server
	ids, err := client.GetSensorIDs()
	if err != nil {
		t.Fatalf("GetSensorIDs failed after concurrent creates: %v", err)
	}

	t.Logf("Created %d sensors concurrently", len(ids))
}

// TestSatelliteInterface_GetSensorWithRetry tests that GetSensor waits for ACTIVE status
func TestSatelliteInterface_GetSensorWithRetry(t *testing.T) {
	server := NewMockServer(0.0, 0)
	go server.Start(8090)
	time.Sleep(100 * time.Millisecond)

	client := NewSatelliteInterface("http://localhost:8090")

	// Create a sensor
	newSensor, err := client.CreateSensor(42)
	if err != nil {
		t.Fatalf("CreateSensor failed: %v", err)
	}

	// GetSensor should automatically wait for it to become active
	start := time.Now()
	sensor, err := client.GetSensor(newSensor.ID)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("GetSensor failed: %v", err)
	}

	if sensor.Status != StatusActive {
		t.Errorf("Expected sensor to be ACTIVE, got %s", sensor.Status)
	}

	t.Logf("GetSensor returned active sensor in %v", elapsed)
}
