package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing activity name.")
			fmt.Println("Usage: tick start <activity_name>")
			return
		}
		activityName := strings.Join(os.Args[2:], " ")
		startActivity(activityName)

	case "end", "stop":
		endActivity()

	case "status":
		showCurrentActivity()

	case "watch":
		watchActivity()

	case "history", "log", "logs":
		showHistory()

	case "stats", "stat":
		period := ""
		if len(os.Args) >= 3 {
			period = os.Args[2]
		}
		showStats(period)

	case "help", "--help", "-h":
		showHelp()

	default:
		fmt.Printf("Unknown command: \"%s\"\n\n", os.Args[1])
		showHelp()
	}
}

func showHelp() {
	fmt.Println("Tick - Simple & Fast CLI Time Tracker")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  tick <command> [arguments]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  start <activity>   Start tracking an activity (only one can run at a time)")
	fmt.Println("  end, stop          End the active activity and record it to history")
	fmt.Println("  status             Show the current running activity and elapsed time")
	fmt.Println("  watch              Live terminal progress display (Ctrl+C quits cleanly)")
	fmt.Println("  history            Show all tracked activities in chronological order")
	fmt.Println("  stats [period]     Show activity sum and total time spent")
	fmt.Println("                     Flags:")
	fmt.Println("                       -d, --day     Today's statistics (default)")
	fmt.Println("                       -w, --week    This week's statistics")
	fmt.Println("                       -m, --month   This month's statistics")
	fmt.Println("                       -a, --all     All-time statistics")
	fmt.Println("  help, -h, --help   Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  tick start coding")
	fmt.Println("  tick status")
	fmt.Println("  tick watch")
	fmt.Println("  tick end")
	fmt.Println("  tick history")
	fmt.Println("  tick stats -d")
	fmt.Println("  tick stats -w")
	fmt.Println("  tick stats -m")
	fmt.Println("  tick stats -a")
}