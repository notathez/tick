package main

import (
	"fmt"
)

// showHistory lists all recorded activities in chronological order.
func showHistory() {
	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	// Filter out invalid/empty entries
	var validEntries []HistoryEntry
	for _, entry := range tickData.History {
		if entry.Activity != "" && !entry.Start.IsZero() {
			validEntries = append(validEntries, entry)
		}
	}

	if len(validEntries) == 0 {
		fmt.Println("No activity history found.")
		return
	}

	var totalDuration int64 = 0

	fmt.Println("==========================================================================================")
	fmt.Println("                                   ACTIVITY HISTORY                                       ")
	fmt.Println("==========================================================================================")
	fmt.Printf(" %-4s %-20s %-12s %-10s %-10s %-15s\n", "#", "ACTIVITY", "DATE", "START", "END", "DURATION")
	fmt.Println("------------------------------------------------------------------------------------------")

	for i, entry := range validEntries {
		totalDuration += entry.Duration
		fmt.Printf(" %-4d %-20s %-12s %-10s %-10s %-15s\n",
			i+1,
			truncateString(entry.Activity, 20),
			entry.Start.Format("2006-01-02"),
			entry.Start.Format("15:04:05"),
			entry.End.Format("15:04:05"),
			formatDuration(entry.Duration),
		)
	}

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Printf(" Total Entries: %d | Total Logged Time: %s\n", len(validEntries), formatDuration(totalDuration))
	fmt.Println("==========================================================================================")
}
