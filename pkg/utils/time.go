package utils

import (
	"source-base/pkg/constraints"
	"time"
)

// ToDuration converts any integer type (seconds) to time.Duration.
func ToDuration[T constraints.Integer](seconds T) time.Duration {
	return time.Duration(seconds) * time.Second
}

// ToDurationMs converts any integer type (milliseconds) to time.Duration.
func ToDurationMs[T constraints.Integer](milliseconds T) time.Duration {
	return time.Duration(milliseconds) * time.Millisecond
}
