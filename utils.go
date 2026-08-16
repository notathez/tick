package main

import (
	"fmt"
	"strings"
)

// formatDuration formats seconds into a human-readable string (e.g. 1d 2h 3m 4s, 2h 15m 30s, 4m 20s, or 35s).
func formatDuration(totalSeconds int64) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}

	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %02dm %02ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %02ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// formatDurationClock formats seconds into a digital timer format (HH:MM:SS).
func formatDurationClock(totalSeconds int64) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// renderProgressBar generates an ASCII progress bar for visualizing percentages.
func renderProgressBar(percentage float64, width int) string {
	if width <= 0 {
		width = 20
	}
	filledLen := int((percentage / 100.0) * float64(width))
	if filledLen > width {
		filledLen = width
	}
	if filledLen < 0 {
		filledLen = 0
	}
	emptyLen := width - filledLen

	return strings.Repeat("█", filledLen) + strings.Repeat("░", emptyLen)
}

// truncateString cuts off a string if it exceeds maxLen and appends "..."
func truncateString(s string, maxLen int) string {
	if len(s) > maxLen && maxLen > 3 {
		return s[:maxLen-3] + "..."
	}
	return s
}
