package main

import (
	"fmt"
	"sort"
)

// DataPoint represents a single point in a time series
type DataPoint struct {
	Time  int
	Value int
}

// TimeSeries represents a series of data points
type TimeSeries []DataPoint

// TaggedDataPoint represents a data point with its series ID
type TaggedDataPoint struct {
	Time     int
	Value    int
	SeriesID int
}

// MergeTimeSeries merges multiple time series by summing values at each unique timestamp
// For each timestamp, it sums the most recent value from each series
// Original implementation using index tracking
func MergeTimeSeries(series ...TimeSeries) TimeSeries {
	if len(series) == 0 {
		return TimeSeries{}
	}

	// 1. Collect all unique timestamps
	timestampSet := make(map[int]bool)
	for _, s := range series {
		for _, point := range s {
			timestampSet[point.Time] = true
		}
	}

	// 2. Convert to sorted slice
	timestamps := make([]int, 0, len(timestampSet))
	for t := range timestampSet {
		timestamps = append(timestamps, t)
	}
	sort.Ints(timestamps)

	// 3. For each timestamp, sum the active values from all series
	result := make(TimeSeries, 0, len(timestamps))

	// Track the current index for each series
	seriesIndices := make([]int, len(series))

	for _, timestamp := range timestamps {
		sum := 0

		// For each series, find the most recent value at this timestamp
		for i, s := range series {
			// Advance index to the latest point <= current timestamp
			for seriesIndices[i] < len(s)-1 && s[seriesIndices[i]+1].Time <= timestamp {
				seriesIndices[i]++
			}

			// If this series has data at or before this timestamp, add its value
			if seriesIndices[i] < len(s) && s[seriesIndices[i]].Time <= timestamp {
				sum += s[seriesIndices[i]].Value
			}
		}

		result = append(result, DataPoint{Time: timestamp, Value: sum})
	}

	return result
}

// MergeTimeSeriesFlatten merges time series using a flatten-and-sort approach
// 1. Flatten all series into one big list with series IDs
// 2. Sort by timestamp
// 3. Process sequentially, tracking last value for each series
// 4. Output sum at each unique timestamp
func MergeTimeSeriesFlatten(series ...TimeSeries) TimeSeries {
	if len(series) == 0 {
		return TimeSeries{}
	}

	// 1. Flatten all data points into one list with series tags
	var allPoints []TaggedDataPoint
	for seriesID, s := range series {
		for _, point := range s {
			allPoints = append(allPoints, TaggedDataPoint{
				Time:     point.Time,
				Value:    point.Value,
				SeriesID: seriesID,
			})
		}
	}

	// 2. Sort by timestamp
	sort.Slice(allPoints, func(i, j int) bool {
		return allPoints[i].Time < allPoints[j].Time
	})

	if len(allPoints) == 0 {
		return TimeSeries{}
	}

	// 3. Track the last known value for each series
	lastValue := make(map[int]int) // seriesID -> last value
	result := make(TimeSeries, 0)

	// Process points, collecting all updates at the same timestamp
	i := 0
	for i < len(allPoints) {
		currentTime := allPoints[i].Time

		// Update all values that occur at this timestamp
		for i < len(allPoints) && allPoints[i].Time == currentTime {
			lastValue[allPoints[i].SeriesID] = allPoints[i].Value
			i++
		}

		// Sum up all current values
		sum := 0
		for _, value := range lastValue {
			sum += value
		}

		result = append(result, DataPoint{Time: currentTime, Value: sum})
	}

	return result
}

// Helper function to print time series
func PrintTimeSeries(name string, series TimeSeries) {
	fmt.Printf("%s: ", name)
	for i, point := range series {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("(%d, %d)", point.Time, point.Value)
	}
	fmt.Println()
}

func main() {
	// Example from the problem
	series1 := TimeSeries{
		{Time: 10, Value: 10},
		{Time: 20, Value: 30},
	}

	series2 := TimeSeries{
		{Time: 15, Value: 20},
	}

	fmt.Println("=== Comparing Both Implementations ===")
	fmt.Println()
	PrintTimeSeries("Series 1", series1)
	PrintTimeSeries("Series 2", series2)

	merged1 := MergeTimeSeries(series1, series2)
	PrintTimeSeries("Merged (Original) ", merged1)

	merged2 := MergeTimeSeriesFlatten(series1, series2)
	PrintTimeSeries("Merged (Flatten)  ", merged2)

	fmt.Println("\n--- Additional Test Cases ---")

	// Test case 2: Three series
	s1 := TimeSeries{
		{Time: 0, Value: 5},
		{Time: 10, Value: 15},
	}

	s2 := TimeSeries{
		{Time: 5, Value: 10},
		{Time: 15, Value: 20},
	}

	s3 := TimeSeries{
		{Time: 8, Value: 3},
		{Time: 12, Value: 7},
	}

	fmt.Println("\nTest Case 2: Three series")
	PrintTimeSeries("Series 1", s1)
	PrintTimeSeries("Series 2", s2)
	PrintTimeSeries("Series 3", s3)
	merged3 := MergeTimeSeries(s1, s2, s3)
	PrintTimeSeries("Merged (Original) ", merged3)
	merged4 := MergeTimeSeriesFlatten(s1, s2, s3)
	PrintTimeSeries("Merged (Flatten)  ", merged4)

	// Test case 3: Overlapping timestamps
	s4 := TimeSeries{
		{Time: 1, Value: 100},
		{Time: 5, Value: 200},
	}

	s5 := TimeSeries{
		{Time: 1, Value: 50},
		{Time: 5, Value: 75},
	}

	fmt.Println("\nTest Case 3: Overlapping timestamps")
	PrintTimeSeries("Series 1", s4)
	PrintTimeSeries("Series 2", s5)
	merged5 := MergeTimeSeries(s4, s5)
	PrintTimeSeries("Merged (Original) ", merged5)
	merged6 := MergeTimeSeriesFlatten(s4, s5)
	PrintTimeSeries("Merged (Flatten)  ", merged6)

	// Test case 4: Single series
	s6 := TimeSeries{
		{Time: 1, Value: 10},
		{Time: 2, Value: 20},
		{Time: 3, Value: 30},
	}

	fmt.Println("\nTest Case 4: Single series (should be unchanged)")
	PrintTimeSeries("Series 1", s6)
	merged7 := MergeTimeSeries(s6)
	PrintTimeSeries("Merged (Original) ", merged7)
	merged8 := MergeTimeSeriesFlatten(s6)
	PrintTimeSeries("Merged (Flatten)  ", merged8)
}
