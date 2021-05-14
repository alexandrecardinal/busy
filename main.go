// Package busy measures the computing's availability of the current host.
// It does so by sleeping and capturing the timestamp at each millisecond. In theory,
// the measurements should show one timestamp per millisecond. In practice, if the host is
// busy, we can estimate the host's scheduling gaps.
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	t0 := time.Now()
	var resolution int64 = 1
	timestamps := getTimestamps(1000*3, resolution)
	gaps := findGaps(timestamps, resolution)
	fmt.Printf("Found %v gaps\n", len(gaps))

	var strGaps string = ""
	jsonGaps, _ := json.Marshal(gaps)
	strGaps = string(jsonGaps)
	fmt.Println(strGaps)

	t1 := time.Now()
	fmt.Printf("Ran in %.2f seconds\n", (t1.Sub(t0)).Seconds())
}

// Returns a list of timestamps for the given duration with a specified resolution (sleep time)
// duration: the duration of the test in milliseconds
// resolution: the frequency at which measurements are taken
func getTimestamps(duration int64, resolution int64) []int64 {
	var ns_duration int64 = duration * 1000 * 1000
	var ns_resolution int64 = resolution * 1000 * 1000
	var iterations int64 = ns_duration / ns_resolution
	var measurements = make([]int64, iterations)
	var current_timestamp int64 = makeTimestamp()

	fmt.Printf("Iterations: %v\n", iterations)
	fmt.Printf("Sleep time: %v ms\n", resolution)
	var i int64
	for i = 0; i < iterations; i++ {
		current_timestamp = makeTimestamp()
		measurements[i] = current_timestamp
		time.Sleep(time.Duration(ns_resolution))
	}
	return measurements
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func findGaps(timestamps []int64, resolution int64) []int64 {
	ns_resolution := resolution * 1000 * 1000
	fmt.Printf("Evaluating %v timestamps\n", len(timestamps))
	var gaps = make([]int64, 0, len(timestamps))
	for index, current_timestamp := range timestamps {
		if index < len(timestamps)-1 {
			next_timestamp := timestamps[index+1]
			diff := next_timestamp - current_timestamp
			if diff > ns_resolution {
				gaps = append(gaps, diff)
			}
		}
	}
	return gaps
}
