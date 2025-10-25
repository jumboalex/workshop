#!/bin/bash

# Script to collect n unique measurements from a sensor
# Usage: ./get-measurements.sh <sensor-id> <n>

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <sensor-id> <n>"
    echo "  sensor-id: The ID of the sensor to query"
    echo "  n:         Number of unique measurements to collect"
    echo ""
    echo "Example: $0 1 5"
    exit 1
fi

SENSOR_ID=$1
N=$2

if [ "$N" -lt 1 ]; then
    echo "Error: n must be at least 1"
    exit 1
fi

echo "Collecting $N unique measurements from sensor $SENSOR_ID..."

declare -A measurements
count=0
attempts=0
max_attempts=$((N * 10))

while [ $count -lt $N ] && [ $attempts -lt $max_attempts ]; do
    attempts=$((attempts + 1))

    # Run get-sensor and extract the measurement value
    output=$(./satellite get-sensor --id=$SENSOR_ID 2>&1)

    # Check if command failed
    if [ $? -ne 0 ]; then
        echo "Error: Failed to get sensor data"
        echo "$output"
        exit 1
    fi

    # Extract measurement (handles both "Measurement: X.XXX" and "Measurement: null")
    measurement=$(echo "$output" | grep "Measurement:" | awk '{print $2}')

    if [ "$measurement" != "null" ] && [ -n "$measurement" ]; then
        # Check if this is a new unique measurement
        if [ -z "${measurements[$measurement]}" ]; then
            measurements[$measurement]=1
            count=$((count + 1))
            echo "  [$count/$N] $measurement"
        fi
    fi

    # Sleep if we haven't collected all measurements yet
    if [ $count -lt $N ]; then
        sleep 1
    fi
done

if [ $count -lt $N ]; then
    echo ""
    echo "Warning: Only collected $count unique measurements after $attempts attempts"
    exit 1
else
    echo ""
    echo "Successfully collected $N unique measurements"
fi
