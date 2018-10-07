package util

import (
	"time"
)

// Now get the time string of time RFC3339
func Now() string {
	return time.Now().Format(time.RFC3339)
}

// UTCNow get the utc time string of time RFC3339
func UTCNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
