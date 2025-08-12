package shared

import (
	"time"
)

// Time represents a timestamp stored as epoch seconds
type Time struct {
	value int64 // epoch seconds
}

// NewTime creates a new Time from a Go time.Time
func NewTime(t time.Time) Time {
	return Time{value: t.Unix()}
}

// NewTimeFromEpoch creates a new Time from epoch seconds
func NewTimeFromEpoch(epoch int64) Time {
	return Time{value: epoch}
}

// NewTimeFromEpochMilli creates a new Time from epoch milliseconds
func NewTimeFromEpochMilli(epochMilli int64) Time {
	return Time{value: epochMilli / 1000}
}

// Now creates a new Time with the current timestamp
func Now() Time {
	return Time{value: time.Now().Unix()}
}

// ToEpoch returns the epoch seconds representation
func (t Time) ToEpoch() int64 {
	return t.value
}

// ToEpochMilli returns the epoch milliseconds representation
func (t Time) ToEpochMilli() int64 {
	return t.value * 1000
}

// ToTime returns the underlying time.Time value
func (t Time) ToTime() time.Time {
	return time.Unix(t.value, 0)
}

// String returns the RFC3339 formatted string representation
func (t Time) String() string {
	return time.Unix(t.value, 0).Format(time.RFC3339)
}

// IsZero reports whether t represents the zero time instant
func (t Time) IsZero() bool {
	return t.value == 0
}

// Before reports whether the time instant t is before u
func (t Time) Before(u Time) bool {
	return t.value < u.value
}

// After reports whether the time instant t is after u
func (t Time) After(u Time) bool {
	return t.value > u.value
}

// Equal reports whether t and u represent the same time instant
func (t Time) Equal(u Time) bool {
	return t.value == u.value
}

// Add returns the time t+d
func (t Time) Add(d time.Duration) Time {
	return Time{value: t.value + int64(d.Seconds())}
}

// Sub returns the duration t-u
func (t Time) Sub(u Time) time.Duration {
	return time.Duration(t.value-u.value) * time.Second
}
