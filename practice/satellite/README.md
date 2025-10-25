# Satellite Sensor API Client

A Go client for managing satellite sensors via an unreliable JSON REST API.

## Features

- Automatic retry logic with exponential backoff
- Handles unreliable/slow connections
- Waits for sensors to reach ACTIVE status automatically
- Simple CLI for managing sensors

## Building

```bash
go build -o satellite .
```

## Usage

### Start Mock Server

```bash
./satellite server --port=8080 --unreliability=0.2 --slowness=500ms
```

### Create a Sensor

```bash
./satellite create-sensor --frequency=42 --url=http://localhost:8080
```

Output:
```
Sensor created successfully!
  ID: 1
  Frequency: 42
  Status: INITIALIZING
```

### Get Sensor Details

Note: This will automatically wait for the sensor to become ACTIVE.

```bash
./satellite get-sensor --id=1 --url=http://localhost:8080
```

Output:
```
Sensor details:
  ID: 1
  Frequency: 42
  Status: ACTIVE
  Measurement: 123.456
```

### List All Sensors

```bash
./satellite list-sensors --url=http://localhost:8080
```

Output:
```
Found 2 sensor(s):
  - 1
  - 2
```

### Run Demo

Starts a mock server and runs a client demo:

```bash
./satellite demo --port=8080 --unreliability=0.2
```

## API

The satellite API provides three endpoints:

- `GET /sensor-ids` - Returns list of sensor IDs
- `POST /sensors` - Creates a sensor with specified frequency
- `GET /sensors/<id>` - Returns sensor details

Sensor states: `INITIALIZING`, `ACTIVE`, `FAILED`, `RESTARTING`, `TERMINATING`

## Testing

```bash
go test -v
```

## Shell Script Helper

### Collect Multiple Measurements

Use the `get-measurements.sh` script to collect n unique measurements from a sensor:

```bash
./get-measurements.sh <sensor-id> <n>
```

**Example:**
```bash
./get-measurements.sh 1 5
```

**Output:**
```
Collecting 5 unique measurements from sensor 1...
  [1/5] 78.348
  [2/5] 77.847
  [3/5] 75.457
  [4/5] 78.704
  [5/5] 77.084

Successfully collected 5 unique measurements
```

The script:
- Calls `satellite get-sensor` repeatedly until it collects n unique measurements
- Tracks unique values using a hash map
- Waits 1 second between requests
- Has a safety limit of n×10 attempts to prevent infinite loops
