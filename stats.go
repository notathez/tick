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
// If activityFilter is non-empty, it delegates to showActivityStats for a focused single-activity view.
func showStats(periodArg string, activityFilter string) {
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

	if activityFilter != "" {
		showActivityStats(tickData, activityFilter, periodTitle, startTime, endTime)
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

// showActivityStats displays daily aggregated time spent on a specific activity.
func showActivityStats(tickData TickData, activityFilter string, periodTitle string, startTime time.Time, endTime time.Time) {
	filter := strings.ToLower(strings.TrimSpace(activityFilter))

	type dayStat struct {
		date     string    // YYYY-MM-DD
		dayTime  time.Time // for formatting & sorting
		duration int64
		sessions int
		hasActive bool
	}

	dailyMap := make(map[string]*dayStat)
	var totalDuration int64
	var totalSessions int

	for _, entry := range tickData.History {
		if entry.Activity == "" || entry.Start.IsZero() {
			continue
		}
		if !strings.Contains(strings.ToLower(entry.Activity), filter) {
			continue
		}
		if (entry.Start.Equal(startTime) || entry.Start.After(startTime)) && entry.Start.Before(endTime) {
			dayKey := entry.Start.Format("2006-01-02")
			stat, exists := dailyMap[dayKey]
			if !exists {
				stat = &dayStat{
					date:    dayKey,
					dayTime: entry.Start,
				}
				dailyMap[dayKey] = stat
			}
			stat.duration += entry.Duration
			stat.sessions++
			totalDuration += entry.Duration
			totalSessions++
		}
	}

	// Include currently active session if it matches
	activeActivity := tickData.Current.Activity
	var activeElapsed int64
	if activeActivity != "" && strings.Contains(strings.ToLower(activeActivity), filter) {
		if (tickData.Current.Start.Equal(startTime) || tickData.Current.Start.After(startTime)) && tickData.Current.Start.Before(endTime) {
			activeElapsed = int64(time.Since(tickData.Current.Start).Seconds())
			if activeElapsed > 0 {
				dayKey := tickData.Current.Start.Format("2006-01-02")
				stat, exists := dailyMap[dayKey]
				if !exists {
					stat = &dayStat{
						date:    dayKey,
						dayTime: tickData.Current.Start,
					}
					dailyMap[dayKey] = stat
				}
				stat.duration += activeElapsed
				stat.sessions++
				stat.hasActive = true
				totalDuration += activeElapsed
				totalSessions++
			}
		}
	}

	fmt.Println("==========================================================================================")
	fmt.Printf(" ACTIVITY STATS: \"%s\" | %s\n", activityFilter, periodTitle)
	fmt.Println("==========================================================================================")

	if len(dailyMap) == 0 {
		fmt.Println(" No activity tracked for this during the selected period.")
		fmt.Println("------------------------------------------------------------------------------------------")
		fmt.Println(" Tip: Activity name is matched by partial, case-insensitive search.")
		fmt.Println("==========================================================================================")
		return
	}

	var daysList []*dayStat
	for _, stat := range dailyMap {
		daysList = append(daysList, stat)
	}

	// Sort days chronologically
	sort.Slice(daysList, func(i, j int) bool {
		return daysList[i].date < daysList[j].date
	})

	fmt.Printf(" %-12s %-10s %-15s %-10s %-25s\n", "DATE", "DAY", "TIME SPENT", "SESSIONS", "SHARE")
	fmt.Println("------------------------------------------------------------------------------------------")

	for _, d := range daysList {
		pct := 0.0
		if totalDuration > 0 {
			pct = (float64(d.duration) / float64(totalDuration)) * 100.0
		}
		timeStr := formatDuration(d.duration)
		if d.hasActive {
			timeStr += " *"
		}
		fmt.Printf(" %-12s %-10s %-15s %-10d %6.1f%%   %-20s\n",
			d.date,
			d.dayTime.Format("Monday"),
			timeStr,
			d.sessions,
			pct,
			renderProgressBar(pct, 15),
		)
	}

	fmt.Println("------------------------------------------------------------------------------------------")

	avgPerDay := totalDuration / int64(len(daysList))
	avgPerSession := int64(0)
	if totalSessions > 0 {
		avgPerSession = totalDuration / int64(totalSessions)
	}

	fmt.Printf(" Total Time: %s across %d day(s) (%d session(s))\n",
		formatDuration(totalDuration),
		len(daysList),
		totalSessions,
	)
	fmt.Printf(" Daily Average: %s/active day | Avg per session: %s\n",
		formatDuration(avgPerDay),
		formatDuration(avgPerSession),
	)
	if activeElapsed > 0 {
		fmt.Printf(" (* Includes %s currently in progress)\n", formatDuration(activeElapsed))
	}
	fmt.Println("==========================================================================================")
}
