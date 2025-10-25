# Quick Start Guide - Satellite Sensor API Client

## 🚀 Get Started in 60 Seconds

### 1. Build the Project
```bash
go build -o satellite .
```

### 2. Run the Demo
```bash
./satellite -mode=demo -unreliability=0.2 -slowness=500ms
```

### 3. Basic Code Usage
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Create client
    client := NewSatelliteClient(ClientConfig{
        BaseURL:    "http://your-satellite-url.com",
        Timeout:    60 * time.Second,
        MaxRetries: 5,
        RetryDelay: 3 * time.Second,
    })

    ctx := context.Background()

    // Create sensor
    sensor, _ := client.CreateSensor(ctx, 42)
    fmt.Printf("Created sensor %d\n", sensor.ID)

    // Get sensor
    s, _ := client.GetSensor(ctx, sensor.ID)
    fmt.Printf("Status: %s\n", s.Status)
}
```

## 📚 Common Operations

### Create a Sensor
```go
sensor, err := client.CreateSensor(ctx, frequency)
```

### Get All Sensor IDs
```go
ids, err := client.GetSensorIDs(ctx)
```

### Get Sensor Details
```go
sensor, err := client.GetSensor(ctx, sensorID)
```

### Wait for Sensor to Activate
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

activeSensor, err := client.WaitForSensorActive(ctx, sensorID, 5*time.Second)
```

### Get All Sensors (Bulk)
```go
sensors, err := client.GetAllSensors(ctx)
for _, s := range sensors {
    fmt.Printf("Sensor %d: %s\n", s.ID, s.Status)
}
```

## ⚙️ Configuration Options

```go
ClientConfig{
    BaseURL:    string,        // Required: API base URL
    Timeout:    time.Duration, // Default: 30s
    MaxRetries: int,           // Default: 5
    RetryDelay: time.Duration, // Default: 2s
}
```

## 🧪 Testing

### Run All Tests
```bash
go test -v
```

### Run with Race Detection
```bash
go test -v -race
```

### Run Specific Test
```bash
go test -v -run TestSatelliteClient_CreateSensor
```

## 🎯 CLI Modes

### Demo Mode (Both Server + Client)
```bash
./satellite -mode=demo -unreliability=0.3 -slowness=2s
```

### Server Only
```bash
./satellite -mode=server -port=8080 -unreliability=0.2 -slowness=1s
```

### Client Only
```bash
./satellite -mode=client -url=http://localhost:8080
```

## ⚠️ Error Handling

### Check for Specific Errors
```go
sensor, err := client.GetSensor(ctx, 999)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // Handle 404
    } else if strings.Contains(err.Error(), "max retries exceeded") {
        // Handle retry exhaustion
    } else if err == context.DeadlineExceeded {
        // Handle timeout
    }
}
```

### Best Practice: Always Use Context with Timeout
```go
// ✓ Good
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
sensor, err := client.CreateSensor(ctx, 42)

// ✗ Bad (could hang forever)
sensor, err := client.CreateSensor(context.Background(), 42)
```

## 🔍 Debugging

### Enable Verbose Logging
The client automatically logs retry attempts:
```
Attempt 1 failed: connection refused
Retry attempt 1/5 for POST /sensors
Retry attempt 2/5 for POST /sensors
✓ Success on attempt 3
```

### Monitor Sensor Status Changes
```go
ticker := time.NewTicker(5 * time.Second)
for range ticker.C {
    s, _ := client.GetSensor(ctx, sensorID)
    fmt.Printf("Status: %s\n", s.Status)
}
```

## 💡 Tips & Tricks

### 1. Handling Nullable Measurements
```go
if sensor.Measurement != nil {
    fmt.Printf("Measurement: %.3f\n", *sensor.Measurement)
} else {
    fmt.Println("No measurement yet")
}
```

### 2. Graceful Degradation
```go
sensors, _ := client.GetAllSensors(ctx)
// Returns successfully retrieved sensors even if some failed
```

### 3. Optimal Retry Configuration
For unreliable satellite (30% failure rate):
```go
ClientConfig{
    Timeout:    60 * time.Second,  // Slow connection
    MaxRetries: 5,                  // More retries
    RetryDelay: 3 * time.Second,   // Longer backoff
}
```

For stable connection:
```go
ClientConfig{
    Timeout:    10 * time.Second,  // Fast fail
    MaxRetries: 2,                  // Fewer retries
    RetryDelay: 1 * time.Second,   // Shorter backoff
}
```

## 📊 Sensor Status Lifecycle

```
INITIALIZING → ACTIVE → (measurements available)
            ↘ FAILED
            ↘ RESTARTING → ACTIVE
                        ↘ TERMINATING
```

## 🐛 Common Issues

### Issue: "max retries exceeded"
**Solution:** Increase `MaxRetries` or check server availability

### Issue: Context timeout
**Solution:** Increase timeout or check network latency
```go
ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
```

### Issue: Sensor stuck in INITIALIZING
**Solution:** Use `WaitForSensorActive()` with appropriate timeout
```go
activeSensor, err := client.WaitForSensorActive(ctx, id, 5*time.Second)
```

## 📖 Further Reading

- See `README.md` for detailed documentation
- See `SUMMARY.md` for architecture decisions
- See `examples.go` for 7 complete usage examples
