package utils

import "time"

// FormattedTimeNow returns the current date and time formatted as "2006-01-02-15:04:05".
// This format is useful for timestamping files or logs.
func FormattedTimeNow() string {
	now := time.Now()
	formattedTime := now.Format("2006-01-02-15:04:05")

	return formattedTime
}
