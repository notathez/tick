package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ActivityStat holds aggregate statistics for a specific activity.
type ActivityStat struct {
	Activity   string
	Duration   int64
	Count      int
	Percentage float64
}

// showStats calculates and displays aggregated time tracking metrics by day, week, month, or all-time.
func showStats(periodArg string) {
	tickData, err := load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	periodArg = strings.ToLower(strings.TrimSpace(periodArg))
	now := time.Now()

	var periodTitle string
	var startTime, endTime time.Time

	switch periodArg {
	case "-w", "--week", "week", "w":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday as 7th day of the week
		}
		startTime = time.Date(now.Year(), now.Month(), now.Day()-(weekday-1), 0, 0, 0, 0, now.Location())
		endTime = startTime.AddDate(0, 0, 7)
		periodTitle = fmt.Sprintf("This Week (%s - %s)", startTime.Format("Jan 02"), startTime.AddDate(0, 0, 6).Format("Jan 02, 2006"))

	case "-m", "--month", "month", "m":
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endTime = startTime.AddDate(0, 1, 0)
		periodTitle = fmt.Sprintf("This Month (%s)", now.Format("January 2006"))

	case "-a", "--all", "all", "a":
		startTime = time.Time{}
		endTime = now.AddDate(100, 0, 0)
		periodTitle = "All Time"

	case "-d", "--day", "day", "today", "d", "":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = startTime.AddDate(0, 0, 1)
		periodTitle = fmt.Sprintf("Today (%s)", now.Format("2006-01-02"))

	default:
		fmt.Printf("Unknown stats option: '%s'\n\n", periodArg)
		fmt.Println("Available options:")
		fmt.Println("  -d, --day     Today's statistics (default)")
		fmt.Println("  -w, --week    This week's statistics")
		fmt.Println("  -m, --month   This month's statistics")
		fmt.Println("  -a, --all     All-time statistics")
		return
	}

	activityDurations := make(map[string]int64)
	activityCounts := make(map[string]int)
	var totalDuration int64 = 0

	for _, entry := range tickData.History {
		if entry.Activity == "" || entry.Start.IsZero() {
			continue
		}

		if (entry.Start.Equal(startTime) || entry.Start.After(startTime)) && entry.Start.Before(endTime) {
			activityDurations[entry.Activity] += entry.Duration
			activityCounts[entry.Activity]++
			totalDuration += entry.Duration
		}
	}

	// Include current active activity if running within the requested period
	var activeElapsed int64 = 0
	activeActivity := tickData.Current.Activity
	if activeActivity != "" {
		if (tickData.Current.Start.Equal(startTime) || tickData.Current.Start.After(startTime)) && tickData.Current.Start.Before(endTime) {
			activeElapsed = int64(time.Since(tickData.Current.Start).Seconds())
			if activeElapsed > 0 {
				activityDurations[activeActivity] += activeElapsed
				activityCounts[activeActivity]++
				totalDuration += activeElapsed
			}
		}
	}

	if totalDuration == 0 && len(activityDurations) == 0 {
		fmt.Println("==========================================================================================")
		fmt.Printf(" STATISTICS: %s\n", periodTitle)
		fmt.Println("==========================================================================================")
		fmt.Println(" No activity tracked during this period.")
		fmt.Println("------------------------------------------------------------------------------------------")
		fmt.Println(" Options: -d (day), -w (week), -m (month), -a (all)")
		fmt.Println("==========================================================================================")
		return
	}

	var statsList []ActivityStat
	for act, dur := range activityDurations {
		pct := 0.0
		if totalDuration > 0 {
			pct = (float64(dur) / float64(totalDuration)) * 100.0
		}
		statsList = append(statsList, ActivityStat{
			Activity:   act,
			Duration:   dur,
			Count:      activityCounts[act],
			Percentage: pct,
		})
	}

	// Sort activities by duration descending
	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].Duration > statsList[j].Duration
	})

	fmt.Println("==========================================================================================")
	fmt.Printf(" STATISTICS: %s\n", periodTitle)
	fmt.Println("==========================================================================================")
	fmt.Printf(" %-20s %-15s %-8s %-8s %-25s\n", "ACTIVITY", "DURATION", "SESSIONS", "SHARE", "DISTRIBUTION")
	fmt.Println("------------------------------------------------------------------------------------------")

	for _, stat := range statsList {
		displayName := stat.Activity
		if displayName == activeActivity && activeElapsed > 0 {
			displayName += " *"
		}

		fmt.Printf(" %-20s %-15s %-8d %6.1f%%   %-25s\n",
			truncateString(displayName, 20),
			formatDuration(stat.Duration),
			stat.Count,
			stat.Percentage,
			renderProgressBar(stat.Percentage, 20),
		)
	}

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Printf(" TOTAL TIME SPENT: %s across %d activities\n", formatDuration(totalDuration), len(statsList))
	if activeActivity != "" && activeElapsed > 0 {
		fmt.Printf(" (* Includes %s currently in progress for \"%s\")\n", formatDuration(activeElapsed), activeActivity)
	}
	fmt.Println("==========================================================================================")
}
