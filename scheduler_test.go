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