package main

// schedule schedules a set of new requests based on a sliding window rate limiting algorithm.
// TODO: Support other algorithms (?)
//
// Parameters:
//
// numCalls (int): The number of new requests to schedule.
// requestCount (int): The maximum number of requests allowed within the specified time frame.
// timeFrame (int): The duration in seconds of the sliding time window.
// previousCalls ([]int64): A slice of Unix timestamps (in seconds) representing the previous requests.
// now (int64): The current Unix timestamp (in seconds).
//
// Returns:
//
// delays ([]int): A slice of delays (in seconds) for each new request.
// previousCalls ([]int64): The updated slice of previous requests, including the new ones.func schedule(numCalls, requestCount, timeFrame int, previousCalls []int64, now int64) ([]int, []int64) {
func schedule(numCalls, requestCount, timeFrame int, previousCalls []int64, now int64) ([]int, []int64) {
	var delays []int

	// Prune previous calls to only keep those within the current time frame
	previousCalls = filterRecentCalls(previousCalls, now-int64(timeFrame))

	// Calculate available slots
	currentCalls := len(previousCalls)
	availableSlots := requestCount - currentCalls

	// Schedule new calls
	for i := 0; i < numCalls; i++ {
		if availableSlots > 0 {
			// No delay if slots are available
			delays = append(delays, 0)
			previousCalls = append(previousCalls, now)
			availableSlots--
		} else {
			// Calculate delay for the next available slot
			nextAvailableTime := previousCalls[0] + int64(timeFrame)
			delay := int(nextAvailableTime - now)
			delays = append(delays, delay)

			// Shift the oldest call to the new scheduled time
			previousCalls = append(previousCalls[1:], nextAvailableTime)
		}
	}

	return delays, previousCalls
}

// filterRecentCalls filters timestamps to keep only those within the current time frame.
func filterRecentCalls(calls []int64, start int64) []int64 {
	var filtered []int64
	for _, t := range calls {
		if t > start {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
