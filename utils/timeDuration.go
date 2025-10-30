package utils

import "fmt"

func IntToTimeDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	return fmt.Sprintf("%2dh %2dm", hours, minutes)
}

func IntToTimeDurationWithSeconds(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60

	return fmt.Sprintf("%2dh %2dm %2ds", hours, minutes, seconds)
}
