package main

import (
	"reflect"
	"testing"
)

func TestSchedule(t *testing.T) {
	tests := []struct {
		numCalls     int
		requestCount int
		timeFrame    int
		expected     []int
	}{
		{
			numCalls:     5,
			requestCount: 2,
			timeFrame:    1,
			expected:     []int{0, 0, 1, 1, 2},
		},
		{
			numCalls:     300,
			requestCount: 100,
			timeFrame:    60,
			expected:     append(append(makeDelaySlice(100, 0), makeDelaySlice(100, 60)...), makeDelaySlice(100, 120)...),
		},
		{
			numCalls:     250,
			requestCount: 100,
			timeFrame:    30,
			expected:     append(append(makeDelaySlice(100, 0), makeDelaySlice(100, 30)...), makeDelaySlice(50, 60)...),
		},
	}

	for _, tt := range tests {
		delays, _ := schedule(tt.numCalls, tt.requestCount, tt.timeFrame, []int64{}, 1729954499)
		if !reflect.DeepEqual(delays, tt.expected) {
			t.Errorf("calculateDelays(%d, %d, %d) = %v; want %v", tt.numCalls, tt.requestCount, tt.timeFrame, delays, tt.expected)
		}
	}
}

// makeDelaySlice creates a slice of size `n` where each element is set to `delay`
func makeDelaySlice(n, delay int) []int {
	slice := make([]int, n)
	for i := range slice {
		slice[i] = delay
	}
	return slice
}


func TestSequentialScheduling(t *testing.T) {
	// Initial parameters (3 calls per minute rate limit)
	requestCount := 3
	timeFrame := 60
	numCalls := 2

	// Initial state for "previous calls"
	previousCalls := []int64{}
	now := int64(0) // Starting time

	// First scheduling call
	expectedFirst := []int{0, 0} // Both calls can happen immediately
	delays, updatedCalls := schedule(numCalls, requestCount, timeFrame, previousCalls, now)
	if !reflect.DeepEqual(delays, expectedFirst) {
		t.Errorf("First schedule call = %v; want %v", delays, expectedFirst)
	}

	// Wait for 30 seconds (simulate passage of time) before the next scheduling
	now += 30

	// Second scheduling call
	expectedSecond := []int{0, 30} // One call immediately, one delayed by 30s due to limited slots
	delays, updatedCalls = schedule(numCalls, requestCount, timeFrame, updatedCalls, now)
	if !reflect.DeepEqual(delays, expectedSecond) {
		t.Errorf("Second schedule call = %v; want %v", delays, expectedSecond)
	}

	// Wait for another 30 seconds (start a new minute)
	now += 30

	// Third scheduling call
	expectedThird := []int{0, 30} // There was 2 calls in the past 60 seconds, we can schedule one and the other will be delayed by 30s
	delays, updatedCalls = schedule(numCalls, requestCount, timeFrame, updatedCalls, now)
	if !reflect.DeepEqual(delays, expectedThird) {
		t.Errorf("Third schedule call = %v; want %v", delays, expectedThird)
	}
}
