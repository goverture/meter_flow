package main

func schedule(numCalls, requestCount, timeFrame int) []int {
	batchCount := numCalls / requestCount
	remainingCalls := numCalls % requestCount
	var delays []int

	// Schedule full batches with a delay after each batch
	for i := 0; i < batchCount; i++ {
		for j := 0; j < requestCount; j++ {
			delays = append(delays, i*timeFrame) // Each call in the batch gets the same delay
		}
	}

	// Schedule any remaining calls with one additional delay
	if remainingCalls > 0 {
		finalDelay := batchCount * timeFrame
		for k := 0; k < remainingCalls; k++ {
			delays = append(delays, finalDelay)
		}
	}

	return delays
}
