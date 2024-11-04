package model

type Resource struct {
	Name           string
	RequestCount   int     // Maximum requests allowed
	TimeFrame      int     // Time frame in seconds
	ScheduledCalls []int64 // Track scheduled timestamps for this resource
}
