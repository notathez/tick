package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// startActivity starts tracking a new activity.
// It ensures that no other activity is currently running.
func startActivity(activity string) {
	if activity == "" {
		fmt.Println("Error: Activity name cannot be empty.")
		fmt.Println("Usage: tick start <activity_name>")
		return
	}

	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	if tickData.Current.Activity != "" {
		elapsed := int64(time.Since(tickData.Current.Start).Seconds())
		fmt.Println("Cannot start a new activity while another is in progress.")
		fmt.Printf("Current active activity: '%s' (running for %s)\n", tickData.Current.Activity, formatDuration(elapsed))
		fmt.Println("Use 'tick end' or 'tick stop' to finish it first.")
		return
	}

	now := time.Now()
	tickData.Current = CurrentActivity{
		Activity: activity,
		Start:    now,
	}

	if err := save(tickData); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
		return
	}

	fmt.Printf("Started activity: \"%s\" at %s\n", activity, now.Format("15:04:05"))
}

// endActivity stops the currently running activity and saves it to history.
func endActivity() {
	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	if tickData.Current.Activity == "" {
		fmt.Println("No activity is currently running.")
		fmt.Println("Use 'tick start <activity_name>' to start one.")
		return
	}

	now := time.Now()
	durationSecs := int64(now.Sub(tickData.Current.Start).Seconds())
	if durationSecs < 0 {
		durationSecs = 0
	}

	entry := HistoryEntry{
		Activity: tickData.Current.Activity,
		Start:    tickData.Current.Start,
		End:      now,
		Duration: durationSecs,
	}

	tickData.History = append(tickData.History, entry)
	endedActivity := tickData.Current.Activity
	tickData.Current = CurrentActivity{}

	if err := save(tickData); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
		return
	}

	fmt.Printf("Ended activity: \"%s\"\n", endedActivity)
	fmt.Printf("  Started:  %s\n", entry.Start.Format("15:04:05"))
	fmt.Printf("  Ended:    %s\n", entry.End.Format("15:04:05"))
	fmt.Printf("  Duration: %s\n", formatDuration(durationSecs))
}

// showCurrentActivity prints the current activity and elapsed time.
func showCurrentActivity() {
	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	if tickData.Current.Activity == "" {
		fmt.Println("No activity is currently running.")
		fmt.Println("Use 'tick start <activity_name>' to start tracking.")
		return
	}

	now := time.Now()
	elapsed := int64(now.Sub(tickData.Current.Start).Seconds())

	fmt.Println("Current Activity:")
	fmt.Printf("  Activity: %s\n", tickData.Current.Activity)
	fmt.Printf("  Started:  %s (%s)\n", tickData.Current.Start.Format("15:04:05"), tickData.Current.Start.Format("2006-01-02"))
	fmt.Printf("  Elapsed:  %s\n", formatDuration(elapsed))
}

// watchActivity displays a real-time progress screen for the current activity.
// Pressing Ctrl+C exits watch mode while keeping the activity running in the background.
func watchActivity() {
	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	if tickData.Current.Activity == "" {
		fmt.Println("No activity is currently running to watch.")
		fmt.Println("Use 'tick start <activity_name>' to start tracking.")
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Initial render
	renderWatchDisplay(tickData.Current)

	for {
		select {
		case <-sigChan:
			fmt.Print("\r\n\n")
			fmt.Printf("Exited watch mode. Activity \"%s\" is still running in the background.\n", tickData.Current.Activity)
			fmt.Println("Use 'tick status' to check time or 'tick end' to finish it.")
			return
		case <-ticker.C:
			renderWatchDisplay(tickData.Current)
		}
	}
}

func renderWatchDisplay(current CurrentActivity) {
	elapsed := int64(time.Since(current.Start).Seconds())
	if elapsed < 0 {
		elapsed = 0
	}

	fmt.Print("\033[H\033[2J")
	fmt.Println("==================================================")
	fmt.Println("             TICK - LIVE WATCH                    ")
	fmt.Println("==================================================")
	fmt.Printf(" Activity : %s\n", current.Activity)
	fmt.Printf(" Started  : %s (%s)\n", current.Start.Format("15:04:05"), current.Start.Format("2006-01-02"))
	fmt.Printf(" Elapsed  : %s  [%s]\n", formatDuration(elapsed), formatDurationClock(elapsed))
	fmt.Println("==================================================")
	fmt.Println(" Press [Ctrl+C] to exit watch mode")
	fmt.Println(" (Activity will keep running in the background)")
}