# Satellite Sensor API Client - Usage Examples

## Basic Usage

```go
// Create a new client
client := NewSatelliteInterface("http://localhost:8080")

// Create a sensor
sensor, err := client.CreateSensor(42)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created sensor %d\n", sensor.ID)

// Wait for sensor to become active (simple retry logic)
activeSensor, err := client.WaitForSensorActive(
    sensor.ID,
    30*time.Second,  // max wait time
    3*time.Second,   // poll interval
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Sensor is now ACTIVE with measurement: %.3f\n", *activeSensor.Measurement)
```

## Manual Polling (if you need custom logic)

```go
// Get sensor details
sensor, err := client.GetSensor(sensorID)
if err != nil {
    log.Fatal(err)
}

// Check status manually
if sensor.Status == StatusActive {
    fmt.Println("Sensor is active!")
}
```

## List All Sensors

```go
ids, err := client.GetSensorIDs()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d sensors: %v\n", len(ids), ids)
```
