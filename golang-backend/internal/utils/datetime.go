package utils

import (
	"math/rand"
	"time"
)

func CurrentMSTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetRandomMSTimestamp(start, end time.Time) int64 {
	if start.After(end) {
		start, end = end, start
	}
	startTimestamp := start.UnixNano() / int64(time.Millisecond)
	endTimestamp := end.UnixNano() / int64(time.Millisecond)

	diff := endTimestamp - startTimestamp

	randomOffset := rand.Int63n(diff)

	randomTimestamp := startTimestamp + randomOffset

	return randomTimestamp
}
