package main

import "time"

// CurrentActivity represents an activity currently being tracked.
type CurrentActivity struct {
	Activity string    `json:"activity"`
	Start    time.Time `json:"start"`
}

// HistoryEntry represents a finished activity entry.
type HistoryEntry struct {
	Activity string    `json:"activity"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration int64     `json:"duration"` // duration in seconds
}

// TickData represents the root structure stored in data.json.
type TickData struct {
	Current CurrentActivity `json:"current"`
	History []HistoryEntry  `json:"history"`
}